<template>
  <a-card>
    <a-spin tip="Loading..." :spinning="pushing">
      <a-tabs default-active-key="1" @change="changeTab">
        <a-tab-pane key="1" force-render>
          <span slot="tab">
            <a-icon type="home" />执行命令
          </span>
          <a-row :gutter="8" v-if="selectKey==='1'">
            <a-col :span="6" style="padding-right: 5px">
              <span>选择DB实例：</span> 
              <a-select v-model="redisDB" style="width: 120px">
                <a-select-option v-for="v in redisDBList" :key="v" :value="v">{{ v }}</a-select-option>
              </a-select>
              <div style="margin-top: 5px;height: 450px;overflow: scroll;border:1px solid #a5b6c8;border-left-width: 0; border-right-width: 1px">
                <s-tree
                  showIcon
                  :dataSource="orgTree"
                  @click="handleClick">
                </s-tree>
              </div>
            </a-col>
            <a-col :span="18">
              <codemirror ref="myCm" :options="cmOptions"></codemirror>
            </a-col>
          </a-row>
          <a-row>
            <a-col :span="24">
              <a-auto-complete @select="onSelect" @search="handleSearch">
                <template slot="dataSource">
                  <a-select-option disabled key="null" v-if="dataSource.length > 0">可执行以下命令</a-select-option>
                  <a-select-option v-for="v in dataSource" :key="v">{{ v }}</a-select-option>
                </template>
                <a-textarea
                  style="margin: 5px;height: 100px;width: 1000px;"
                  v-model="redis_cmd"
                  @keydown.tab.native="handleKeyTab"
                  placeholder="输入执行命令..." />
              </a-auto-complete>
            </a-col>
            <a-col :span="24">
              <a-button style="left: 10px;" type="primary" size="large" @click="handleCmd">执行</a-button>
            </a-col>
          </a-row>
        </a-tab-pane>
        <a-tab-pane key="2" tab="实例信息">
          <a-row :gutter="8" v-if="selectKey==='2'">
            <a-col :span="6" style="padding-right: 5px">
              <div style="margin-top: 5px;height: 650px;overflow: scroll;border:1px solid #a5b6c8;border-left-width: 0; border-right-width: 1px">
                <s-tree
                  showIcon
                  :dataSource="orgTree"
                  @click="getRedisInfo">
                </s-tree>
              </div>
            </a-col>
            <a-col :span="18">
              <a-descriptions bordered size="small" :column="2" >
                <a-descriptions-item v-for="(v, k) in redisMetrics" :key="k" :value="v" :label="k">
                  <span>{{ v }}</span>
                </a-descriptions-item>
              </a-descriptions>
            </a-col>
          </a-row>
        </a-tab-pane>
        <a-tab-pane key="3" tab="健康检测">
          <a-row :gutter="8" v-if="selectKey==='3'">
            <a-col :span="6" style="padding-right: 5px">
              <div style="margin-top: 5px;height: 650px;overflow: scroll;border:1px solid #a5b6c8;border-left-width: 0; border-right-width: 1px">
                <s-tree showIcon
                  :dataSource="orgTree"
                  @click="handleClick">
                </s-tree>
              </div>
            </a-col>
            <a-col :span="18">
              <a-card title="检测进度">
                  <a-button style="left: 10px;" type="primary" :disabled="disButton" @click="checkHealth">开始检测</a-button>
                  <a-progress style="right: 10px;" :percent="progress" />
              </a-card>
              <a-card title="检测结果" style="margin-top: 10px;">
                <div v-for="(item, key) in redisHealth" :key="key">
                  <span v-if="item.length>0">{{ item }}</span>
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </a-spin>
  </a-card>
</template>

<script>
import STree from '@/components/Tree/Tree'
import redisApi from "@/api/redisms.js"

import 'codemirror/theme/material.css'
import "codemirror/mode/javascript/javascript"
import "codemirror/addon/selection/active-line"
import "codemirror/addon/display/autorefresh"
import "codemirror/keymap/sublime"
import notification from 'ant-design-vue/es/notification'



export default {
  name: 'RedisList',
  components: {
    STree,
  },
  data () {
    return {
      orgTree: [],
      result: [],
      openKeys: null,
      selectKey: '1',
      redis_cmd: null,
      redisDB: 0,
      redisDBList: [...new Array(16).keys()],
      redisMetrics: null,
      redisHealth: ["暂无数据"],
      pushing: false,
      progress: 0,
      disButton: false,
      redisCmds: [],
      dataSource: [],
      cmOptions: {
        mode: "text/x-shell",
        readOnly: true,
        autoRefresh: true,
        autofocus: false,
        keyMap: "sublime", // 编辑器模式
        autoCloseBrackets: true,
        lineWrapping: true, // 代码折叠
        foldGutter: true,
        theme: "material",
        gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"]
      },
      checkOptions: ["cluster_status", "durability", "bulk_ops", "metrics"],
    }
  },
  created () {
    this.fetchRedisIns()
    this.getRedisCmdList()
  },
  mounted() {
    this.$refs.myCm.codemirror.setValue("查询结果输出 ...")
  },
  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror
    },
  },
  methods: {
    fetchRedisIns() {
      redisApi.getRedisIns().then(resp => {
        this.orgTree = resp.data
      }).catch(error => {
        notification.error({
          message: 'error',
          description: "无法加载redis列表"
        })
      })
    },
    getRedisCmdList() {
      redisApi.getRedisCmdList().then(resp => {
        this.redisCmds = resp.data
      })
    },
    changeTab(e) {
      this.openKeys = null
      this.selectKey = e
    },
    handleClick (e) {
      this.openKeys = e.key
      //console.log(this.$refs.myCm.codemirror)
    },
    getRedisInfo(e) {
      this.redisMetrics = null
      this.pushing = true
      redisApi.getRedisCmd(e.key, this.redisDB).then(resp => {
        this.redisMetrics = resp.data
      }).catch(error => {
        notification.error({
          message: 'error',
          description: "无法加载redis Metrics"
        })
      }).finally(() => {
        this.pushing = false
      })
    },
    onSelect(value) {
      this.redis_cmd = value
    },
    handleSearch(value) {
      this.dataSource = []
      if(value && value.length >1){
        for(var i in this.redisCmds) {
          if(this.redisCmds[i].indexOf(value)>=0){
            this.dataSource.push(this.redisCmds[i])
          }
        }
      }
    },
    handleKeyTab(e) {
      e.preventDefault()
      if(this.dataSource.length > 0) {
        e.target.value = this.dataSource[0]
      }
    },
    handleCmd() {
      if(this.openKeys === null){
        this.$message.error("请先选择redis实例")
        return false
      }

      if(this.redis_cmd === null){
        this.$message.error("无效的执行命令")
        return false
      }
      const data = {
        "redis_id": this.openKeys,
        "redis_cmd": this.redis_cmd,
        "redis_db": this.redisDB,
      }
      this.pushing = true
      redisApi.execRedisCmd(data).then(resp => {
        if(resp.code === "0000") {
          const data = resp.data
          for(var item of data["data"]) {
            const formatResult = data["redis_name"]+ "> " +item["redis_cmd"]+ "\n" +this.formatData(item["result"], true)+"\n"
            this.result.push(formatResult)
          }
          this.$refs.myCm.codemirror.setValue(this.result.join(""))
          this.$refs.myCm.codemirror.scrollTo(0, this.$refs.myCm.codemirror.doc.height) //移动右侧滚动条
        } else {
          this.$message.error(resp.message)
        }
      }).catch(error => {
        this.$message.error("内部错误")
      }).finally(() => {
        this.pushing = false
      })
    },
    checkHealth() {
      if(this.openKeys === null){
        this.$message.error("请先选择redis实例")
        return false
      } else {
        this.checkDetailHealth()
      }
    },
    async checkDetailHealth() {
      this.redisHealth = []
      this.progress = 0
      this.disButton = true 
      for(let option of this.checkOptions) {
        await redisApi.getRedisHealth(this.openKeys, option).then(resp => {
          this.progress += Math.floor(Math.random() * 10) + 5
          if(resp.code === "0000") {
            if(resp.data.length>0) {
              this.redisHealth.push(resp.data)
            }
          } else {
            this.$message.error(resp.message)
          }
        }).catch(error => {
          this.$message.error("内部错误")
        })
      }
      if(this.redisHealth.length===0) {
        this.redisHealth = ["检测完成，暂未发现问题"]
      }
      this.progress = 100
      this.disButton = false  
    },
    formatData(obj, flag) {
      if (!obj && typeof(obj) != "undefined" && obj !== 0) {
        return "None"
      } else if ((typeof(obj) == "object") && obj.constructor === Array) {
        let obj_str = ""
        if(obj.length === 0){
          obj_str = "[]\n" 
        } else {
          for (const [k, value] of Object.entries(obj)) {
            const key = parseInt(k)
            if (flag) {
              obj_str += (key + 1) + ") " + this.formatData(value, false) + "\n"
            } else if (key === 0) {
              obj_str += (key + 1) + ") " + this.formatData(value) + "\n"
            } else {
              obj_str += "   "+ (key + 1) + ") " + this.formatData(value) + "\n"
            }
          }
        }
        return obj_str
      } else if ((typeof(obj) == "object") && obj.constructor === Object) {
        let obj_str = ""
        if (Object.keys(obj).length === 0) {
          obj_str = "{}\n"
        } else {
          for (const [key, value] of Object.entries(obj)) {
            if (flag) {
              obj_str += "  "+ key + ": " + this.formatData(value, false) + "\n"
            } else {
              obj_str += key + "=" + this.formatData(value, false) + ","
            }
          }
        }
        return obj_str
      } else {
        return obj
      }
    }
  }
}
</script>

<style lang="less">
  .custom-tree {

    /deep/ .ant-menu-item-group-title {
      position: relative;
      &:hover {
        .btn {
          display: block;
        }
      }
    }

    /deep/ .ant-menu-item {
      font-size: 12px;
      &:hover {
        .btn {
          display: block;
        }
      }
    }

    /deep/ .btn {
      display: none;
      position: absolute;
      top: 0;
      right: 10px;
      width: 20px;
      height: 40px;
      line-height: 40px;
      z-index: 1050;

      &:hover {
        transform: scale(1.2);
        transition: 0.5s all;
      }
    }
  }
  
  .CodeMirror {
    border: 1px solid #eee;
    height: 487px;
    font-size: 13px;
  }

  .spin-content {
    border: 1px solid #91d5ff;
    background-color: #e6f7ff;
    padding: 30px;
  }

  .ant-form-item-control {
    height: 10;
    line-height: 10px;
  }

  .ant-select-search__field__wrap { 
    overflow: hidden;
  }

</style>