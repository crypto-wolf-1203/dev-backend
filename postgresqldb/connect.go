package psqldb

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"pongpongi.com/osdep"
	"pongpongi.com/blockchain/evm"
	"fmt"
)

type dbQueryType map[string]interface{}
var db *sqlx.DB

func InitDB() {
	var err error

	dbURL := osdep.Getenv("DATABASE_URL", "postgres://devdb:devdb@localhost/backenddb?sslmode=disable")
	fmt.Println(dbURL)
	db, err = sqlx.Open("postgres", dbURL)
	osdep.Check(err)

	{ // history table
		query := fmt.Sprintf("create table if not exists %s (ip text not null, method text not null, content text not null, timestamp timestamptz primary key not null default now());",
				GetHistoryTableNable())
		DbQuery(query)
	}

	{ // evm chain table
		query := fmt.Sprintf("create table if not exists evm_chain_table (name text not null primary key, nativecurrency text not null, chainid integer, rpcurl text not null, blockexplorer text not null);")
		DbQuery(query)

		registerChain("bsc")
		registerChain("bsctestnet")
	}

	fmt.Println("PostgreSQL successfully connected!")
}

func GetHistoryTableNable() string {
	return "connection_history"
}

func registerChain(name string) {
	bscInfo, err := evm.GetChainInfo(name)
	if err == nil {
		query := fmt.Sprintf("select * from evm_chain_table where name='%s';", bscInfo.Name)
		retChains := DbQuery(query)

		if len(retChains) == 0 {
			query = fmt.Sprintf("insert into evm_chain_table(name, nativecurrency, rpcurl, chainid, blockexplorer) values('%s', '%s', '%s', %d, '%s');", bscInfo.Name, bscInfo.NativeCurrency, bscInfo.RpcURL, bscInfo.ChainId, bscInfo.BlockExplorer)
			DbQuery(query)
		} else {
			query = fmt.Sprintf("update evm_chain_table set nativecurrency = '%s', rpcurl = '%s', chainid=%d, blockexplorer = '%s' where name='%s';", bscInfo.NativeCurrency, bscInfo.RpcURL, bscInfo.ChainId, bscInfo.BlockExplorer, bscInfo.Name)
			DbQuery(query)
		}
	}
}
