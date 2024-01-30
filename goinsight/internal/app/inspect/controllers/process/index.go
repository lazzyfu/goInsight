/*
@Time    :   2022/07/06 10:12:48
@Author  :   zongfei.fu
@Desc    :   None
*/

package process

import (
	"fmt"
	"goInsight/global"
	"goInsight/internal/app/inspect/config"
	"goInsight/internal/pkg/kv"
	"goInsight/internal/pkg/utils"
	"strings"

	"github.com/jinzhu/copier"
)

// 检查索引前缀
type IndexPrefix struct {
	Table         string
	UniqueKeys    []string
	SecondaryKeys []string
	FulltextKeys  []string
	InspectParams *config.InspectParams
}

func (i *IndexPrefix) CheckUniquePrefix() error {
	var unMatchKeys []string
	for _, key := range i.UniqueKeys {
		if i.InspectParams.CHECK_IDENTIFIER {
			if ok := utils.IsMatchPattern(utils.NamePattern, key); !ok {
				return fmt.Errorf("索引`%s`命名不符合要求，仅允许匹配正则`%s`[表`%s`]", key, utils.NamePattern, i.Table)
			}
		}
		if len(key) == 0 {
			return fmt.Errorf("表`%s`必须显式指定唯一索引名称", i.Table)
		}
		if !utils.HasPrefix(key, i.InspectParams.UNQI_INDEX_PREFIX, false) {
			unMatchKeys = append(unMatchKeys, key)
		}
	}
	if len(unMatchKeys) > 0 {
		return fmt.Errorf("唯一索引前缀不符合要求，必须以`%s`开头（不区分大小写）[表`%s`]", i.InspectParams.UNQI_INDEX_PREFIX, i.Table)
	}
	return nil
}

func (i *IndexPrefix) CheckSecondaryPrefix() error {
	var unMatchKeys []string
	for _, key := range i.SecondaryKeys {
		if i.InspectParams.CHECK_IDENTIFIER {
			if ok := utils.IsMatchPattern(utils.NamePattern, key); !ok {
				return fmt.Errorf("索引`%s`命名不符合要求，仅允许匹配正则`%s`[表`%s`]", key, utils.NamePattern, i.Table)
			}
		}
		if len(key) == 0 {
			return fmt.Errorf("表`%s`必须显式指定二级索引名称", i.Table)
		}
		if !utils.HasPrefix(key, i.InspectParams.SECONDARY_INDEX_PREFIX, false) {
			unMatchKeys = append(unMatchKeys, key)
		}
	}
	if len(unMatchKeys) > 0 {
		return fmt.Errorf("二级索引前缀不符合要求，必须以`%s`开头（不区分大小写）[表`%s`]", i.InspectParams.SECONDARY_INDEX_PREFIX, i.Table)
	}
	return nil
}

func (i *IndexPrefix) CheckFulltextPrefix() error {
	var unMatchKeys []string
	for _, key := range i.FulltextKeys {
		if i.InspectParams.CHECK_IDENTIFIER {
			if ok := utils.IsMatchPattern(utils.NamePattern, key); !ok {
				return fmt.Errorf("索引`%s`命名不符合要求，仅允许匹配正则`%s`[表`%s`]", key, utils.NamePattern, i.Table)
			}
		}
		if len(key) == 0 {
			return fmt.Errorf("表`%s`必须显式指定全文索引名称", i.Table)
		}
		if !utils.HasPrefix(key, i.InspectParams.FULLTEXT_INDEX_PREFIX, false) {
			unMatchKeys = append(unMatchKeys, key)
		}
	}
	if len(unMatchKeys) > 0 {
		return fmt.Errorf("全文索引前缀不符合要求，必须以`%s`开头（不区分大小写）[表`%s`]", i.InspectParams.FULLTEXT_INDEX_PREFIX, i.Table)
	}
	return nil
}

// 检查索引数量
type IndexLen struct {
	Index string
	Len   int
}
type IndexNumber struct {
	Table         string
	Number        int        // 二级索引的个数
	Keys          []IndexLen // 存储索引名和组成索引列的个数
	InspectParams *config.InspectParams
}

// 最多有N个二级索引,包括唯一索引
func (i *IndexNumber) CheckSecondaryIndexesNum() error {
	if i.Number > i.InspectParams.MAX_INDEX_KEYS {
		return fmt.Errorf("表`%s`最多允许定义%d个二级索引，当前二级索引个数为%d", i.Table, i.InspectParams.MAX_INDEX_KEYS, i.Number)
	}
	return nil
}

// 主键索引列数不能超过指定的个数
func (i *IndexNumber) CheckPrimaryKeyColsNum() error {
	for _, item := range i.Keys {
		if item.Index == "PrimaryKey" {
			// 主键索引列数不能超过指定的个数
			if item.Len > i.InspectParams.PRIMARYKEY_MAX_KEY_PARTS {
				return fmt.Errorf("表`%s`的主键索引`PRIMARY KEY`最多允许组成列数为%d，当前列数为%d", i.Table, i.InspectParams.PRIMARYKEY_MAX_KEY_PARTS, item.Len)
			}
		} else {
			// 二级索引的列数不能超过指定的个数,包括唯一索引
			if item.Len > i.InspectParams.SECONDARY_INDEX_MAX_KEY_PARTS {
				return fmt.Errorf("表`%s`的二级索引`%s`最多允许组成列数为%d，当前列数为%d", i.Table, item.Index, i.InspectParams.SECONDARY_INDEX_MAX_KEY_PARTS, item.Len)
			}
		}
	}
	return nil
}

// 检查冗余索引
// is_drop：drop操作，is_add：add操作，is_meta：原表表结构
type IndexColsMap struct {
	Index string   // 索引
	Tag   string   // 标记,值:is_drop/is_add/is_meta
	Cols  []string // 组成索引的列
}
type RedundantIndex struct {
	Table       string
	Cols        []string       // 列
	Indexes     []string       // 索引名组合
	IndexesCols []IndexColsMap // 索引名和列名组合
}

func (r *RedundantIndex) CheckRepeatCols() error {
	// 索引中的列,不能重复,不区分大小写,建索引时,指定的列必须存在
	// KEY idx_a (col1,col2,col1),
	for _, item := range r.IndexesCols {
		idxColsDefDup := make(map[string]bool)
		for _, col := range item.Cols {
			itemLower := strings.ToLower(col)
			if !utils.IsContain(r.Cols, itemLower) {
				return fmt.Errorf("索引`%s`中的列`%s`不存在[表`%s`]", item.Index, itemLower, r.Table)
			}
			if !idxColsDefDup[itemLower] {
				idxColsDefDup[itemLower] = true
			} else {
				return fmt.Errorf("索引`%s`中的列不能重复[表`%s`]", item.Index, r.Table)
			}
		}
	}
	return nil
}

func (r *RedundantIndex) CheckRepeatColsWithDiffIndexes() error {
	// 查找重复的索引，即索引名不一样，但是定义的列一样，不区分大小写
	// KEY idx_a_b (col1,col2),
	// KEY idx_b (col1,col2),
	// 获取is_add
	var isAddIdxCols []string
	// 获取is_drop
	var isDropIdxCols []string
	// 获取is_meta
	isMetaIdxCols := make(map[string]string)
	for _, item := range r.IndexesCols {
		if item.Tag == "is_add" {
			isAddIdxCols = append(isAddIdxCols, item.Index)
		}
		if item.Tag == "is_drop" {
			isDropIdxCols = append(isDropIdxCols, item.Index)
		}
		if item.Tag == "is_meta" {
			isMetaIdxCols[item.Index] = strings.ToLower(strings.Join(item.Cols, utils.KeyJoinChar))
		}
	}
	// 从is_meta数据中删除is_drop的数据
	for _, item := range isDropIdxCols {
		delete(isMetaIdxCols, item)
	}
	// 为create table语句
	if len(isDropIdxCols) == 0 && len(isAddIdxCols) == 0 {
		// 对map中的键进行分组
		groups := make(map[string][]string)
		for key, value := range isMetaIdxCols {
			groups[value] = append(groups[value], key)
		}
		var result []string
		for _, keys := range groups {
			if len(keys) > 1 {
				result = append(result, strings.Join(keys, "/"))
			}
		}
		if len(result) > 0 {
			return fmt.Errorf("表`%s`发现了重复定义的索引：%s", r.Table, strings.Join(result, ";"))
		}
	} else {
		// 为alter table xxx add / drop index
		for _, item := range r.IndexesCols {
			// 查找重复的索引,即索引名不一样,但是定义的列一样,不区分大小写
			// KEY idx_a_b (col1,col2),
			// KEY idx_b (col1,col2),
			if item.Tag == "is_add" {
				found := false
				var foundKey string
				valueJoin := strings.ToLower(strings.Join(item.Cols, utils.KeyJoinChar))
				for key, value := range isMetaIdxCols {
					if value == valueJoin {
						found = true
						foundKey = key
						break
					}
				}
				if !found {
					isMetaIdxCols[item.Index] = valueJoin
				} else {
					return fmt.Errorf("表`%s`发现了重复定义的索引：%s(%s)，请检查原始表结构", r.Table, foundKey, strings.Join(item.Cols, ", "))
				}
			}
		}
	}
	return nil
}

func (r *RedundantIndex) CheckRedundantColsWithDiffIndexes() error {
	/*
		查找冗余的索引，即索引名不一样，但是定义的列冗余，不区分大小写
		KEY idx_a (col1),
		KEY idx_b (col1,col2),
		KEY idx_c (col1,col2,col3)

		r.IndexesCols数据结构：
		[
			{UNIQ_PID_COMMITID_TOKEN is_add [i_project_id ch_commit_id i_task_type]}
			{UNIQ_PID_COMMITID is_drop []}
			{UNIQ_PID_COMMITID  [i_project_id ch_commit_id]}
			{IDX_STATUS_RETRY_COUNT  [i_status i_retry_count]}
			{IDX_STATUS_UPDATETIME  [i_status d_update_time]}
		]
	*/
	// idxCols := make(map[string]string)

	/*
		解决冗余索引时，当指定drop冗余索引时，允许通过
		表结构已有索引uniq_aa(`a`, `b`)
		添加索引：
			ALTER TABLE `tbl1` ADD UNIQUE uniq_bb(`a`,`b`,`c`),DROP INDEX `idx_aa`;
	*/
	// 获取is_drop
	var isDropIdxCols []string
	// 获取is_meta
	isMetaIdxCols := make(map[string]string)
	isAddIdxCols := make(map[string]string)

	for _, item := range r.IndexesCols {
		if item.Tag == "is_add" {
			isAddIdxCols[item.Index] = strings.ToLower(strings.Join(item.Cols, utils.KeyJoinChar))
		}
		if item.Tag == "is_drop" {
			isDropIdxCols = append(isDropIdxCols, item.Index)
		}
		if item.Tag == "is_meta" {
			isMetaIdxCols[item.Index] = strings.ToLower(strings.Join(item.Cols, utils.KeyJoinChar))
		}
	}
	// 从is_meta数据中删除is_drop的数据
	for _, item := range isDropIdxCols {
		delete(isMetaIdxCols, item)
	}

	// 为create table语句
	if len(isDropIdxCols) == 0 && len(isAddIdxCols) == 0 {
		var result []string
		for k, v := range isMetaIdxCols {
			for k1, v1 := range isMetaIdxCols {
				if v != v1 {
					if strings.HasPrefix(v, v1) && utils.IsSubKey(v, v1) {
						result = append(result, strings.Join([]string{k, k1}, "/"))
					}
				}
			}
		}
		if len(result) > 0 {
			return fmt.Errorf("表`%s`发现了冗余索引，冗余索引的字段组合为%s", r.Table, strings.Join(result, ", "))
		}
	} else {
		// 为alter table xxx add / drop index
		var result []string
		for k, v := range isAddIdxCols {
			for k1, v1 := range isMetaIdxCols {
				if v != v1 {
					if strings.HasPrefix(v, v1) && utils.IsSubKey(v, v1) {
						result = append(result, strings.Join([]string{k, k1}, "/"))
					}
				}
			}
		}
		if len(result) > 0 {
			return fmt.Errorf("表`%s`发现了冗余索引，冗余索引的字段组合为%s", r.Table, strings.Join(result, ", "))
		}
	}
	return nil
}

// BLOB/TEXT类型不能设置为索引
type DisabledIndexes struct {
	Table       string
	Cols        []string       // 列
	IndexesCols []IndexColsMap // 索引名和列名组合
}

func (i *DisabledIndexes) Check() error {
	if len(i.Cols) > 0 {
		for _, item := range i.IndexesCols {
			for _, col := range item.Cols {
				if utils.IsContain(i.Cols, col) {
					return fmt.Errorf("索引名`%s`中的列`%s`不能创建索引[表`%s`]", item.Index, col, i.Table)
				}
			}
		}
	}
	return nil
}

// IndexLargePrefix
type LargePrefixIndexPartSpecification struct {
	Column  string
	Tp      byte
	Elems   []string // Elems is the element list for enum and set type.
	Ilen    int      // key `idx_name`(name(32))中的32
	Flen    int      // 字段长度
	Decimal int      // decimal字段专用,decimal(12,2)中的2
	Charset string   // 列字符集
}
type LargePrefixIndexColsMap struct {
	Name string // 索引
	Keys []LargePrefixIndexPartSpecification
}
type LargePrefix struct {
	Table                    string // 表名
	Charset                  string // 表字符集
	LargePrefixIndexColsMaps []LargePrefixIndexColsMap
}

func (l *LargePrefix) Check(kv *kv.KVCache) error {
	indexMaxLength := 767

	dbVersion := kv.Get("dbVersion").(string)
	versionIns := DbVersion{dbVersion}
	if versionIns.IsTiDB() {
		indexMaxLength = 3072
	} else {
		if versionIns.Int() > 80000 {
			indexMaxLength = 3072
		}
		if versionIns.Int() > 50700 && kv.Get("largePrefix").(string) == "ON" {
			indexMaxLength = 3072
		}
	}
	for _, i := range l.LargePrefixIndexColsMaps {
		// &{meta_cluster utf8 [{idx_datacenter [{datacenter 254 -1 128 } {cluster_domain 15 32 128 }]}]}
		var maxSumLength int
		for _, key := range i.Keys {
			// 判断字符集，当列字符集为空，使用表的字符集
			if len(key.Charset) == 0 {
				key.Charset = l.Charset
			}
			var instDataBytes DataBytes
			err := copier.CopyWithOption(&instDataBytes, key, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			if err != nil {
				return err
			}
			maxSumLength += instDataBytes.Get(versionIns.Int())
		}

		global.App.Log.Debug(fmt.Sprintf("maxSumLength:%d, indexMaxLength:%d", maxSumLength, indexMaxLength))
		if maxSumLength > indexMaxLength {
			return fmt.Errorf("表`%s`的索引`%s`超出了innodb-large-prefix限制，当前索引长度为%d字节，最大限制为%d字节【例如可使用前缀索引，如：Field(length)】", l.Table, i.Name, maxSumLength, indexMaxLength)
		}
	}
	return nil
}
