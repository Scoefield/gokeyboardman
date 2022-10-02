package main

import (
	"comm/ini"
	"comm/md_log"
	"comm/mysqlclient"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Config struct {
	LogLevel   int
	DataSet    string
	DBHost	   string
	TXDBHost   string
	User	   string
	Pwd		   string
}

const (
	DB_POOL_CH = 100
)

var (
	g_pDBPool       *mysqlclient.MysqlDBPool
	g_TXDBPool      *mysqlclient.MysqlDBPool
	g_Config      	Config

	g_StartTime	string
	g_EndTime	string

	g_DealTableInfo *mysqlclient.MysqlResultSet
)

var difCh = make(chan map[string]interface{}, 100)
var dbPoolCh = make(chan map[string]*mysqlclient.MysqlDBPool, DB_POOL_CH)

var tableInfoMap = make(map[string]string)

type Context struct {
	Config struct {
		Uri      string
		Username string
		Password string
		Db       string
	}
	DbConn *sql.DB
}

func CheckErr(err error, msg ...string) {
	if err != nil{
		errMsg := ""
		for _, m:= range msg {
			errMsg += m
		}
		panic(err.Error() + errMsg + "\n")
	}
}

/*
*  比较逻辑
 */
func DoCompareLogic(start, end time.Time, tableName string) {
	startTime, endTime := start.Format(Layout), end.Format(Layout)
	md_log.Keyf("[Starting compare logic] - tableName:%v, startTime:%s, endTime:%s \n", tableName, startTime, endTime)

	dbPoolMap := <-dbPoolCh
	defer func() {
		dbPoolCh <- dbPoolMap
	}()

	dbPool := dbPoolMap["dbPool"]
	txDBPool := dbPoolMap["txDBPool"]

	// 查js cloud db 数据
	dealSql := fmt.Sprintf("select * from %s where Flast_update_time>='%v' and Flast_update_time<'%v'", tableName, startTime, endTime)
	oRs, err := dbPool.Query(dealSql)
	if err != nil {
		md_log.Errorf(-1, nil, "查询deal数据失败, err:%v, tableName:%v, dealSql:%v", err, tableName, dealSql)
		return
	}
	defer oRs.Close()

	var txDealSql string
	var strDealId, pkName string
	primaryIds := make([]int64, 0)

	dealDataMap := make(map[int64]map[string]interface{})
	//dealDataSlice := make([]map[string]interface{}, 0)
	for oRs.Next() {

		dealMap := GetDealDataMap(oRs, tableInfoMap)

		dealId := oRs.GetInt64("Fdeal_id")
		buyerId := oRs.GetInt64("Fbuyer_id")
		strDealId = fmt.Sprintf("%08d%04d%04d", dealId, buyerId%1000, oRs.GetInt64("Fseller_id")%1000)

		var primaryKey int64
		if strings.HasPrefix(tableName, "t_deal") {
			primaryKey = dealId
			primaryIds = append(primaryIds, dealId)
			pkName = "Fdeal_id"
		} else if strings.HasPrefix(tableName, "t_recv") {
			recvId := oRs.GetInt64("Frecv_fee_id")
			primaryKey = recvId
			primaryIds = append(primaryIds, recvId)
			pkName = "Frecv_fee_id"
		} else if strings.HasPrefix(tableName, "t_aftersale") {
			afterSaleId := oRs.GetInt64("Faftersale_id")
			primaryKey = afterSaleId
			primaryIds = append(primaryIds, afterSaleId)
			pkName = "Faftersale_id"
		} else {
			continue
		}
		dealMap["primaryKey"] = primaryKey
		dealMap["strDealId"] = strDealId
		dealDataMap[primaryKey] = dealMap
		//dealDataSlice = append(dealDataSlice, dealDataMap)
	}

	if len(primaryIds) == 0 {
		return
	}
	a, _ := json.Marshal(primaryIds)
	b := strings.ReplaceAll(string(a), "[", "(")
	pIds := strings.ReplaceAll(b, "]", ")")
	txDealSql = fmt.Sprintf("SELECT * FROM %v WHERE %v IN %v ", tableName, pkName, pIds)

	txRet, err := txDBPool.Query(txDealSql)
	if err != nil {
		md_log.Errorf(-1, nil, "查询 txDeal 数据失败, tableName:%v, txDealSql:%v, err:%v", tableName, txDealSql, err)
		return
	}

	txHasRecord := false
	for txRet.Next() {
		CompareField(tableInfoMap, dealDataMap, txRet, pkName, tableName)
		txHasRecord= true
	}
	if !txHasRecord {
		difMap := make(map[string]interface{})
		md_log.Debugf("查询 txDeal 数据 empty, tableName:%v, strDealId:%v", tableName, strDealId)
		difMap["primaryKey"] = string(a)
		difMap["dealId"] = strDealId
		difMap["tableName"] = tableName
		difMap["desc"] = fmt.Sprintf("txCloud no record")
		difCh <- difMap
	}
	txRet.Close()


	md_log.Keyf("[End compare logic] - tableName:%v, startTime:%s, endTime:%s \n", tableName, startTime, endTime)

	return
}

func CompareField(tableInfoMap map[string]string, dealDataMap map[int64]map[string]interface{}, txRet *mysqlclient.MysqlResultSet, pkName, tableName string) {
	difMap := make(map[string]interface{})

	primaryKey := txRet.GetInt64(pkName)
	dealData := dealDataMap[primaryKey]

	difMap["tableName"] = tableName
	difMap["primaryKey"] = primaryKey
	difMap["dealId"] = dealData["strDealId"]


	difFlag := false
	for fieldName, colType := range tableInfoMap {	// 遍历表字段属性
		var a, b interface{}
		//colType := tableInfo.GetString("COLUMN_TYPE")		// 字段类型
		//colType := tableInfo.GetString("DATA_TYPE")		// 数据类型
		//fieldName := tableInfo.GetString("COLUMN_NAME")	// 取字段名

		switch colType {
		case "int", "tinyint":
			a, b = dealData[fieldName].(int32), txRet.GetInt32(fieldName)
		case "bigint":
			a, b = dealData[fieldName].(int64), txRet.GetInt64(fieldName)
		case "varchar", "datetime", "text":
			a, b = dealData[fieldName].(string), txRet.GetString(fieldName)
		default:	// 其它类型非关键字段，可不处理
			continue
		}

		if a != b {
			difMap[fieldName] = fmt.Sprintf("%v:%v|%v", fieldName, a, b)
			difFlag = true
		}
	}
	if difFlag {
		difCh <- difMap
	}
}

// 导出有差异的订单数据
func WriteDiffDeal(wg *sync.WaitGroup, fileName string) {
	defer wg.Done()

	pExport := xlsx.NewFile()
	pSheet, _ := pExport.AddSheet("sheet1")
	headRow := pSheet.AddRow()
	for _, v := range DIFF_DEAL_HEAD {
		outCell := headRow.AddCell()
		outCell.SetValue(v)
	}

	for mapData := range difCh {
		outRow := pSheet.AddRow()
		addCell(outRow, mapData["primaryKey"])
		addCell(outRow, mapData["dealId"])
		addCell(outRow, mapData["tableName"])
		for key, val := range mapData {
			if key == "primaryKey" || key == "dealId" || key == "tableName" {
				continue
			}
			addCell(outRow, val)
		}
	}

	//保存导出文件
	path := fileName + ".xlsx"
	err := pExport.Save(path)
	if err != nil {
		md_log.Errorf(-1, nil, "文件存储失败，拟存储路径 path:[%v]", path)
		return
	}
	md_log.Debugf("WriteDiffDeal done, save file success.")
}

func addCell(outRow *xlsx.Row, n interface{}) {
	outCell := outRow.AddCell()
	outCell.SetValue(n)
}

// 获取表的数据总数
func GetCount(table string) (int, error) {
	cntSql := fmt.Sprintf("select count(*) as cnt from %s where Flast_update_time>='%v' and Flast_update_time<'%v'", table, g_StartTime, g_EndTime)

	cRst, err := g_pDBPool.Query(cntSql)
	if err != nil {
		md_log.Errorf(-1, nil, "查询数量失败, tableName:%v, sql:%v", table, cntSql)
		return 0, err
	}
	defer cRst.Close()

	var cnt int
	for cRst.First() {
		cnt = cRst.GetInt("cnt")
	}

	return cnt, nil
}

// 多协程 处理整张表(WaitGroup 和 chan 组合实现)
// offset, limit 处理方式，效率可能会慢点，不推荐用此方法
func DealByGoroutine(nunRoutine int, tableName string) {
	pageSize := 100
	total, _ := GetCount(tableName)	// err 的话，total 为 0，这里就不再判断和打印 err 了
	loops := total/pageSize + 1		// 分批处理，每批处理 pageSize 条数据

	md_log.Keyf("*** table:%v, total:%v", tableName, total)

	wg := new(sync.WaitGroup)
	ch := make(chan struct{}, nunRoutine)	// 限制最多开 nunRoutine 个 goroutine 处理
	// 这里开 loops 个 goroutine 处理一张表的数据
	for i := 0; i < loops; i++ {
		wg.Add(1)
		ch <- struct{}{}  	// 当协程数达到上限（nunRoutine）时阻塞下面协程的创建，否则继续
		//oft := i * pageSize	// oft: table 偏移, pageSize: table limit
		go func() {
			//err := DoCompareLogic(nil, , tableName)
			//if err != nil {
			//	md_log.Errorf(-1, nil, "[ERROR] tableName:%v, Offset %d, Limit %d, err:%v", tableName, oft, pageSize, err)
			//}
			<- ch
			wg.Done()
		}()
	}
	wg.Wait()
}


// 多协程 处理整张表(WaitGroup 和 chan 组合实现)
// 根据 Flast_update_time 分时间段处理，开多个 goroutine 处理
func DealByGoroutineWithUpTime(tableName string) {
	//md_log.Debugf("DealByGoroutineWithUpTime start. table:%v", tableName)

	parseStartTime, _ := time.Parse(Layout, g_StartTime)
	parseEndTime, _ := time.Parse(Layout, g_EndTime)
	subTime := parseEndTime.Sub(parseStartTime).Hours()

	loops := int(subTime/24) + 1

	wg := new(sync.WaitGroup)
	ch := make(chan struct{}, 10)
	for i := 0; i < loops; i++ {
		wg.Add(1)

		ch <- struct{}{}  	// 当协程数达到上限（nunRoutine）时阻塞下面协程的创建，否则继续
		startTime := parseStartTime.Add(time.Hour * 24 * time.Duration(i))
		endTime := startTime.Add(time.Hour * 24)
		go func(start, end time.Time) {
			defer wg.Done()
			DoCompareLogic(start, end, tableName)
			<- ch
		}(startTime, endTime)
	}
	wg.Wait()
}


// 分组处理表数据
func StepByTable(begin, end int, table string) {
	//0~99, 100~199, ..., 900~999
	for i := begin; i < end; i++ {
		tableName := fmt.Sprintf("%v_%04d", table, i)
		DealByGoroutineWithUpTime(tableName)

		md_log.Debugf("table: [%s] compare done.", tableName)
	}
}

func CompareAfterData(table string) {
	defer func() {
		if e := recover(); e != nil{
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			md_log.Errorf(-1, nil, "CompareDealData panic, errMsg: %v, stackBuf: %v\n", e, string(buf))
		}
	}()

	var wgc sync.WaitGroup

	wgc.Add(1)
	go WriteDiffDeal(&wgc, fmt.Sprintf("dif_%v", table))

	DealByGoroutineWithUpTime(table)

	close(difCh)
	wgc.Wait()

	md_log.Debugf("[CompareDealData All] Done")

}

// 比对订单数据入口
func CompareDealData(table string) {
	defer func() {
		if e := recover(); e != nil{
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			md_log.Errorf(-1, nil, "CompareDealData panic, errMsg: %v, stackBuf: %v\n", e, string(buf))
		}
	}()

	var wg sync.WaitGroup
	var wgc sync.WaitGroup

	wgc.Add(1)
	go WriteDiffDeal(&wgc, fmt.Sprintf("dif_%v", table))

	// *** 用 10 个 goroutine 分别处理 1000 张表（0~999），即每个 goroutine 处理 100 张表 ***
	size := 100					// 每个 goroutine 处理的表数量
	tableTotal := 1000			// 表的数量
	loops := tableTotal / size	// 需要开启的 goroutine 的数量，10
	start, end := 0, 0
	for i := 0; i < loops; i++ {
		// 0~99, 100~199, ..., 900~999
		start = i * size
		end = start + size
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			StepByTable(start, end, table)
		}(start, end)
	}

	wg.Wait()

	close(difCh)
	wgc.Wait()

	for i := 0; i < DB_POOL_CH; i++ {
		dbPoolMap := <-dbPoolCh
		dbpool := dbPoolMap["dbPool"]
		txDBPool := dbPoolMap["txDBPool"]
		dbpool.Close()
		txDBPool.Close()
	}

	md_log.Debugf("[CompareDealData All] Done")
}


func ReadConfigFile(sConfigFile string) error {
	fmt.Println("begin ReadConfigFile, configFile", sConfigFile)
	oFile, err := ini.LoadFile(sConfigFile)
	if err != nil {
		fmt.Println("read config file failed, configFile", sConfigFile)
		return err
	}

	g_Config.LogLevel = oFile.GetInt("base", "logswitch", 255)

	g_Config.DBHost = oFile.Get("biz", "dbHost", "")
	g_Config.TXDBHost = oFile.Get("biz", "txDBHost", "")
	g_Config.User = oFile.Get("biz", "user", "")
	g_Config.Pwd = oFile.Get("biz", "pwd", "")

	fmt.Println("ConfigInit success, g_Config", g_Config)
	return nil
}


func Init(sConfigFile string) error {
	// 读取配置文件
	fmt.Println("usage: ", os.Args[0], os.Args[1])
	err := ReadConfigFile(sConfigFile)
	if err != nil {
		return err
	}

	md_log.Init(DAEMON_NAME, g_Config.LogLevel)
	md_log.Header()
	md_log.Debugf("log init success, LogLevel:%v", g_Config.LogLevel)

	for i := 0; i < DB_POOL_CH; i++ {
		dbPoolMap := make(map[string]*mysqlclient.MysqlDBPool)
		dbPool, err := mysqlclient.NewDBPool(g_Config.User, g_Config.Pwd, g_Config.DBHost, "md_deal", "utf8", 3306, 60)
		if err != nil {
			md_log.Errorf(-1, nil, "mysqlclient.InitDBPoolBySet failed. err:%v", err)
			return err
		}
		txDBPool, err := mysqlclient.NewDBPool(g_Config.User, g_Config.Pwd, g_Config.TXDBHost, "md_deal", "utf8", 3306, 60)
		if err != nil {
			md_log.Errorf(-1, nil, "mysqlclient.InitDBPoolBySet failed. err:%v", err)
			return err
		}
		dbPoolMap["dbPool"] = dbPool
		dbPoolMap["txDBPool"] = txDBPool
		dbPoolCh <- dbPoolMap
	}

	md_log.Debugf("init db conn poll success")
	fmt.Println("init db conn poll success")

	return nil
}

func GetTableInfo(db, tableName string) {
	dbPoolMap := <-dbPoolCh
	defer func() {
		dbPoolCh <- dbPoolMap
	}()

	dbPool := dbPoolMap["dbPool"]

	strSql := fmt.Sprintf("SELECT * FROM INFORMATION_SCHEMA.COLUMNS where TABLE_SCHEMA='%v' and TABLE_NAME='%v'", db, tableName)
	tableInfo, err := dbPool.Query(strSql)
	if err != nil {
		fmt.Println("query INFORMATION_SCHEMA db err", tableName, err)
		return
	}
	defer tableInfo.Close()

	// tableInfoMap := make(map[string]string)
	for tableInfo.Next() {
		colType := tableInfo.GetString("DATA_TYPE")		// 数据类型
		fieldName := tableInfo.GetString("COLUMN_NAME")	// 取字段名
		tableInfoMap[fieldName] = colType
	}
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please usage: ", os.Args[0], " daemon_compare_dealdb_cloud.ini")
		return
	}

	// 初始化相关配置
	sConfigFile := string(os.Args[1])
	err := Init(sConfigFile)
	if err != nil {
		return
	}

	// 处理数据的开始和结束时间
	g_StartTime = os.Args[2]
	g_EndTime = os.Args[3]
	md_log.Debugf("startTime:%v, endTime:%v", g_StartTime, g_EndTime)

	start := time.Now()

	// 比较订单数据入口
	if os.Args[4] == "deal" {
		GetTableInfo("md_deal", "t_deal_0000")
		CompareDealData("t_deal")
	}
	// 售后打款单
	if os.Args[4] == "recv" {
		CompareDealData("t_recv_fee")
	}

	// 比较售后数据入口
	if os.Args[4] == "after" {
		CompareAfterData("t_aftersale")
	}



	fmt.Println("CompareDealData finish. Spend time:", time.Now().Sub(start).Seconds())
	md_log.Keyf("CompareDealData finish. Spend time:%v", time.Now().Sub(start).Seconds())

	return
}
