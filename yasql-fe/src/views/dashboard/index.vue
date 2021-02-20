<template>
  <div>
    <a-card title="系统仪表盘">
      <a-row :gutter="24">
        <a-col :sm="24" :md="12" :xl="14">
          <a-row :gutter="24">
            <a-col :sm="24" :md="12" :xl="12" :style="{ marginBottom: '24px' }">
              <chart-card :loading="loading" title="用户数" :total="sysdata.user_count">
                <template slot="footer"
                  >今日新增：<span>{{ sysdata.user_count_inc }}</span></template
                >
              </chart-card>
            </a-col>

            <a-col :sm="24" :md="12" :xl="12" :style="{ marginBottom: '24px' }">
              <chart-card :loading="loading" title="数据源" :total="sysdata.database_source">
                <template slot="footer"
                  >今日新增：<span>{{ sysdata.database_source_inc }}</span></template
                >
              </chart-card>
            </a-col>

            <a-col :sm="24" :md="12" :xl="12" :style="{ marginBottom: '24px' }">
              <chart-card :loading="loading" title="工单数" :total="sysdata.orders_count">
                <template slot="footer"
                  >今日新增：<span>{{ sysdata.orders_count_inc }}</span></template
                >
              </chart-card>
            </a-col>

            <a-col :sm="24" :md="12" :xl="12" :style="{ marginBottom: '24px' }">
              <chart-card :loading="loading" title="DMS查询次数" :total="sysdata.dms_count">
                <template slot="footer"
                  >今日新增：<span>{{ sysdata.dms_count_inc }}</span></template
                >
              </chart-card>
            </a-col>
          </a-row>
        </a-col>
        <a-col :sm="24" :md="12" :xl="8">
          <div id="sysPieChart"></div>
        </a-col>
      </a-row>
    </a-card>

    <a-card title="我的仪表盘" style="margin-top: 4px">
      <a-row :gutter="24">
        <a-col :sm="24" :md="12" :xl="14">
          <a-row :gutter="24">
           <a-col :sm="24" :md="12" :xl="12" :style="{ marginBottom: '24px' }">
              <chart-card :loading="loading" title="工单数" :total="selfdata.orders_count">
                <template slot="footer"
                  >今日新增：<span>{{ selfdata.orders_count_inc }}</span></template
                >
              </chart-card>
            </a-col>

            <a-col :sm="24" :md="12" :xl="12" :style="{ marginBottom: '24px' }">
              <chart-card :loading="loading" title="DMS查询次数" :total="selfdata.dms_count">
                <template slot="footer"
                  >今日新增：<span>{{ selfdata.dms_count_inc }}</span></template
                >
              </chart-card>
            </a-col>
          </a-row>
        </a-col>
        <a-col :sm="24" :md="12" :xl="8">
          <div id="selfPieChart"></div>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script>
import { ChartCard } from '@/components'
import { getSysDash, getSelfDash } from '@/api/dashboard'

export default {
  components: {
    ChartCard,
  },
  data() {
    return {
      loading: true,
      sysdata: '',
      selfdata: '',
    }
  },
  methods: {
    renderPieCharts(chart, title, data) {
      chart.setOption({
        title: {
          text: title,
          subtext: '',
          left: 'center',
        },
        tooltip: {
          trigger: 'item',
        },
        legend: {
          orient: 'vertical',
          left: 'right',
        },
        series: [
          {
            name: '分类',
            type: 'pie',
            radius: '40%',
            data: data,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)',
              },
            },
          },
        ],
      })
    },
    drawOrderPie() {
      getSysDash().then((response) => {
        if (response.code != '0000') {
          this.$message.error(response.message)
        } else {
          this.sysdata = response.data
          let sysPieChart = this.$echarts.init(document.getElementById('sysPieChart'))
          this.renderPieCharts(sysPieChart, '系统工单占比', this.sysdata.pie_data)
        }

        getSelfDash().then((response) => {
          if (response.code != '0000') {
            this.$message.error(response.message)
          } else {
            this.selfdata = response.data
            let selfPieChart = this.$echarts.init(document.getElementById('selfPieChart'))
            this.renderPieCharts(selfPieChart, '我的工单占比', this.selfdata.pie_data)
          }
        })
      })
    },
  },
  created() {
    setTimeout(() => {
      this.loading = !this.loading
    }, 1000)
  },
  mounted() {
    this.drawOrderPie()
  },
}
</script>

<style scoped>
/* 设置pie宽高 */
#sysPieChart,
#selfPieChart,
#LineChart {
  width: 90%;
  height: 320px;
}
</style>
