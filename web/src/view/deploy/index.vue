<template>
  <div>
    <el-table
      :data="tableData"
      style="width: 100%"
      :cell-style="cellStyle"
    >
      <el-table-column
        prop="deploy_status"
        label="状态"
        width="180"
      />
      <el-table-column
        prop="metadata.name"
        label="名称"
        width="180"
      >
      <template slot-scope="scope">
        <el-button size="mini" type="success" @click="toDetile(scope.row)">{{scope.row.metadata.name}}</el-button>
      </template>
      </el-table-column>
      <el-table-column
        prop="spec.template.spec.containers"
        label="镜像"
      >
        <template slot-scope="{row}">
          <span v-for="(tag,i) in row.spec.template.spec.containers "><el-tag>{{handleStr(tag.image,i)}}</el-tag></span>
        </template>
      </el-table-column>
      <el-table-column
        prop="status.replicas"
        label="副本数"
      >
        <el-table-column
          prop="spec.replicas"
          label="期望"
          width="80"
        />
        <el-table-column
          prop="status.availableReplicas"
          label="available"
          width="80"
        />
        <el-table-column
          prop="status.unavailableReplicas"
          label="unavailable"
          width="100"
        />
      </el-table-column>
    </el-table>
  </div>
</template>

<script>

import { store } from '@/store'
const token = store.getters['user/token']
const user = store.getters['user/userInfo']
const path = process.env.VUE_APP_BASE_API
export default {
  name: 'Index',
  data() {
    return {
      tableData: [],
      replicaInfo: '<span>sdf</span>'
    }
  },
  computed: {
  },
  created() {
    console.log('启动')
    this.getDeployList()
  },
  mounted() {

  },
  methods: {
    handleStr: function(v, i) {
      if (i === 0) {
        return v
      } else if (i === 1) {
        return '+其他'
      } else {
        return ''
      }
    },
    getDeployList() {
      return this.axios.get(path + '/kubernetes/deploy/list', {
        headers: {
          'Content-Type': 'application/json',
          'x-token': token,
          'x-user-id': user.ID
        }
      })
        .then(rsp => {
          this.tableData = rsp.data.items
          for (let i = 0; i < this.tableData.length; i++) {
            if (this.tableData[i].status.availableReplicas === undefined) {
              this.tableData[i].status.availableReplicas = 0
            }
            if (this.tableData[i].status.unavailableReplicas === undefined) {
              this.tableData[i].status.unavailableReplicas = 0
            }
          }
        })
    },
    cellStyle({ row, column, rowIndex, columnIndex }) {
      if (columnIndex === 5) { // 指定坐标rowIndex ：行，columnIndex ：列
        return 'color: red' // rgb(105,0,7)
      } else if (columnIndex === 4) {
        return 'color: green'
      } else {
        return ''
      }
    },
    toDetile(row) {
      this.$router.push({
        name: 'detail',
        params: {
        },
        query: {
          name: row.metadata.name
        }
      })
    }

  }

}
</script>

<style scoped>
.red{
  color: red;
}
.black{
  color: black;
}
</style>
