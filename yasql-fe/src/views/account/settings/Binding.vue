<template>
  <div id="orderPieChart"></div>
</template>

<script>
export default {
  data() {
    return {}
  },
  methods: {
    renderPieCharts(chart) {
      chart.setOption({
        title: {
          text: '某站点用户访问来源',
          subtext: '纯属虚构',
          left: 'center',
        },
        tooltip: {
          trigger: 'item',
        },
        legend: {
          orient: 'vertical',
          left: 'left',
        },
        series: [
          {
            name: '访问来源',
            type: 'pie',
            radius: '50%',
            data: [
              { value: 1048, name: '搜索引擎' },
              { value: 735, name: '直接访问' },
              { value: 580, name: '邮件营销' },
              { value: 484, name: '联盟广告' },
              { value: 300, name: '视频广告' },
            ],
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
      let orderPieChart = this.$echarts.init(document.getElementById('orderPieChart'))
      this.renderPieCharts(orderPieChart)
    },
  },
  mounted() {
    this.drawOrderPie()
  },
}
</script>

<style scoped>

/* 设置pie宽高 */
#orderPieChart,
#tablePieChart,
#LineChart {
  width: 100%;
  height: 350px;
}
</style>
