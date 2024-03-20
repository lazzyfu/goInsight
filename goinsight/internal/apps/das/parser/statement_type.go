/*
@Time    :   2023/04/11 15:21:19
@Author  :   zongfei.fu
@Desc    :
*/

package parser

import (
	"reflect"
	"strings"

	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"
)

type SubExplain struct {
	StatementType string
}

func (s *SubExplain) Enter(in ast.Node) (ast.Node, bool) {
	if stmt, ok := in.(*ast.ExplainStmt); ok {
		switch stmt.Stmt.(type) {
		case *ast.SelectStmt:
			s.StatementType = "ExplainSelect"
		case *ast.DeleteStmt:
			s.StatementType = "ExplainDelete"
		case *ast.InsertStmt:
			s.StatementType = "ExplainInsert"
		case *ast.UpdateStmt:
			s.StatementType = "ExplainUpdate"
		case *ast.SetOprStmt:
			s.StatementType = "ExplainUnion"
		case *ast.ShowStmt:
			s.StatementType = "Desc"
		default:
			s.StatementType = strings.Replace(reflect.TypeOf(stmt.Stmt).String(), "*", "", 1)
		}
	}
	return in, false
}

func (s *SubExplain) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type StatementType struct{}

func (s *StatementType) Extract(in ast.StmtNode) (statementType string) {
	// 判断语句类型是否被允许
	switch stmt := in.(type) {
	case *ast.SelectStmt:
		// WITH CTE AS 属于SELECT
		statementType = "SELECT"
	case *ast.DeleteStmt:
		statementType = "DELETE"
	case *ast.InsertStmt:
		statementType = "INSERT"
	case *ast.UpdateStmt:
		statementType = "UPDATE"
	case *ast.SetOprStmt:
		statementType = "UNION"
	case *ast.ShowStmt:
		if stmt.Tp == ast.ShowEngines {
			statementType = "ShowEngines"
		}
		if stmt.Tp == ast.ShowDatabases {
			statementType = "ShowDatabases"
		}
		if stmt.Tp == ast.ShowTables {
			statementType = "ShowTables"
		}
		if stmt.Tp == ast.ShowTableStatus {
			statementType = "ShowTableStatus"
		}
		if stmt.Tp == ast.ShowColumns {
			statementType = "ShowColumns"
		}
		if stmt.Tp == ast.ShowWarnings {
			statementType = "ShowWarnings"
		}
		if stmt.Tp == ast.ShowCharset {
			statementType = "ShowCharset"
		}
		if stmt.Tp == ast.ShowVariables {
			statementType = "ShowVariables"
		}
		if stmt.Tp == ast.ShowStatus {
			statementType = "ShowStatus"
		}
		if stmt.Tp == ast.ShowCollation {
			statementType = "ShowCollation"
		}
		if stmt.Tp == ast.ShowCreateTable {
			statementType = "ShowCreateTable"
		}
		if stmt.Tp == ast.ShowCreateView {
			statementType = "ShowCreateView"
		}
		if stmt.Tp == ast.ShowCreateUser {
			statementType = "ShowCreateUser"
		}
		if stmt.Tp == ast.ShowCreateSequence {
			statementType = "ShowCreateSequence"
		}
		if stmt.Tp == ast.ShowCreatePlacementPolicy {
			statementType = "ShowCreatePlacementPolicy"
		}
		if stmt.Tp == ast.ShowGrants {
			statementType = "ShowGrants"
		}
		if stmt.Tp == ast.ShowTriggers {
			statementType = "ShowTriggers"
		}
		if stmt.Tp == ast.ShowProcedureStatus {
			statementType = "ShowProcedureStatus"
		}
		if stmt.Tp == ast.ShowIndex {
			statementType = "ShowIndex"
		}
		if stmt.Tp == ast.ShowProcessList {
			statementType = "ShowProcessList"
		}
		if stmt.Tp == ast.ShowCreateDatabase {
			statementType = "ShowCreateDatabase"
		}
		if stmt.Tp == ast.ShowConfig {
			statementType = "ShowConfig"
		}
		if stmt.Tp == ast.ShowEvents {
			statementType = "ShowEvents"
		}
		if stmt.Tp == ast.ShowStatsExtended {
			statementType = "ShowStatsExtended"
		}
		if stmt.Tp == ast.ShowStatsMeta {
			statementType = "ShowStatsMeta"
		}
		if stmt.Tp == ast.ShowStatsHistograms {
			statementType = "ShowStatsHistograms"
		}
		if stmt.Tp == ast.ShowStatsTopN {
			statementType = "ShowStatsTopN"
		}
		if stmt.Tp == ast.ShowStatsBuckets {
			statementType = "ShowStatsBuckets"
		}
		if stmt.Tp == ast.ShowStatsHealthy {
			statementType = "ShowStatsHealthy"
		}
		if stmt.Tp == ast.ShowStatsLocked {
			statementType = "ShowStatsLocked"
		}
		if stmt.Tp == ast.ShowHistogramsInFlight {
			statementType = "ShowHistogramsInFlight"
		}
		if stmt.Tp == ast.ShowColumnStatsUsage {
			statementType = "ShowColumnStatsUsage"
		}
		if stmt.Tp == ast.ShowPlugins {
			statementType = "ShowPlugins"
		}
		if stmt.Tp == ast.ShowProfile {
			statementType = "ShowProfile"
		}
		if stmt.Tp == ast.ShowProfiles {
			statementType = "ShowProfiles"
		}
		if stmt.Tp == ast.ShowMasterStatus {
			statementType = "ShowMasterStatus"
		}
		if stmt.Tp == ast.ShowPrivileges {
			statementType = "ShowPrivileges"
		}
		if stmt.Tp == ast.ShowErrors {
			statementType = "ShowErrors"
		}
		if stmt.Tp == ast.ShowBindings {
			statementType = "ShowBindings"
		}
		if stmt.Tp == ast.ShowBindingCacheStatus {
			statementType = "ShowBindingCacheStatus"
		}
		if stmt.Tp == ast.ShowPumpStatus {
			statementType = "ShowPumpStatus"
		}
		if stmt.Tp == ast.ShowDrainerStatus {
			statementType = "ShowDrainerStatus"
		}
		if stmt.Tp == ast.ShowOpenTables {
			statementType = "ShowOpenTables"
		}
		if stmt.Tp == ast.ShowAnalyzeStatus {
			statementType = "ShowAnalyzeStatus"
		}
		if stmt.Tp == ast.ShowRegions {
			statementType = "ShowRegions"
		}
		if stmt.Tp == ast.ShowBuiltins {
			statementType = "ShowBuiltins"
		}
		if stmt.Tp == ast.ShowTableNextRowId {
			statementType = "ShowTableNextRowId"
		}
		if stmt.Tp == ast.ShowBackups {
			statementType = "ShowBackups"
		}
		if stmt.Tp == ast.ShowRestores {
			statementType = "ShowRestores"
		}
		if stmt.Tp == ast.ShowImports {
			statementType = "ShowImports"
		}
		if stmt.Tp == ast.ShowCreateImport {
			statementType = "ShowCreateImport"
		}
		if stmt.Tp == ast.ShowPlacement {
			statementType = "ShowPlacement"
		}
		if stmt.Tp == ast.ShowPlacementForDatabase {
			statementType = "ShowPlacementForDatabase"
		}
		if stmt.Tp == ast.ShowPlacementForTable {
			statementType = "ShowPlacementForTable"
		}
		if stmt.Tp == ast.ShowPlacementForPartition {
			statementType = "ShowPlacementForPartition"
		}
		if stmt.Tp == ast.ShowPlacementLabels {
			statementType = "ShowPlacementLabels"
		}
		if stmt.Tp == ast.ShowSessionStates {
			statementType = "ShowSessionStates"
		}
	case *ast.ExplainStmt:
		// 遍历explain语句
		v := &SubExplain{}
		(in).Accept(v)
		statementType = v.StatementType
	case *ast.ExplainForStmt:
		statementType = "ExplainFor"
	case *ast.CallStmt:
		statementType = "CALL"
	case *ast.SetStmt:
		statementType = "SET"
	case *ast.AlterTableStmt, *ast.AlterSequenceStmt, *ast.AlterPlacementPolicyStmt:
		statementType = "ALTER"
	case *ast.CreateDatabaseStmt, *ast.CreateIndexStmt, *ast.CreateTableStmt, *ast.CreateViewStmt, *ast.CreateSequenceStmt, *ast.CreatePlacementPolicyStmt:
		statementType = "CREATE"
	case *ast.DropDatabaseStmt, *ast.DropIndexStmt, *ast.DropTableStmt, *ast.DropSequenceStmt, *ast.DropPlacementPolicyStmt:
		statementType = "DROP"
	case *ast.RenameTableStmt:
		statementType = "RENAME"
	case *ast.TruncateTableStmt:
		statementType = "TRUNCATE"
	case *ast.RepairTableStmt:
		statementType = "REPAIR"
	case *ast.LoadDataStmt:
		statementType = "LoadData"
	case *ast.SplitRegionStmt:
		statementType = "SplitRegion"
	case *ast.NonTransactionalDMLStmt:
		statementType = "NonTransactionalDML"
	case *ast.UseStmt:
		statementType = "Use"
	case *ast.ShutdownStmt:
		statementType = "ShutDown"
	default:
		statementType = strings.Replace(reflect.TypeOf(stmt).String(), "*", "", 1)
	}
	return statementType
}
