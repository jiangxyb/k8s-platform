<template>
  <div>
    <p class="bigWord">deploy细节:</p>
    <el-table
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column
        prop="name"
        label="deployment  名称"
        width="180"
      />
      <el-table-column
        prop="images_name"
        label="镜像"
      >
        <template slot-scope="{row}">
          <span v-for="(tag,i) in row.images_name "><el-tag>{{ tag }}</el-tag></span>
        </template>
      </el-table-column>
      <el-table-column
        prop="replicas"
        label="副本数"
      >
        <template slot-scope="tags">
          <span v-for="(tag,i) in tags.row.replicas ">{{ toMy(tag, i, 0) }}</span>
          <span v-for="(tag,i) in tags.row.replicas " class="green">{{ toMy(tag, i, 1) }}</span>
          <span v-for="(tag,i) in tags.row.replicas " class="red">{{ toMy(tag, i, 2) }}</span>
          <el-button size="mini" type="success" @click="incReplica('default',name,false)">+</el-button>
          <el-button size="mini" type="success" @click="incReplica('default',name,true)">-</el-button>
        </template>

      </el-table-column>
      <el-table-column
        prop="replicas"
        label="删除"
      >
        <template slot-scope="t">
          <el-button size="mini" type="success" @click="deleteDeploy(t.row)">Delete</el-button>
        </template>
      </el-table-column>>
    </el-table>
    <div style="height:30px;" />
    <!--   -------------------------------------------------------------------------- -->
    <!--    pod内容-->
    <p class="bigWord">对应Pod:</p>
    <el-table
      :data="pods"
      style="width: 100%"
    >
      <el-table-column
        prop="phase"
        label="阶段"
      >
        <template slot-scope="{row}">
          <p>{{ row.phase }}</p>
          <span class="red">{{ row.message }}</span>
        </template>
      </el-table-column>>

      <el-table-column
        prop="name"
        label="Pod名称"
      >
        <template slot-scope="{row}">
          <p>{{ row.name }}</p>
          <span v-show="!row.ready" class="red">{{ row.event_msg }}</span>
        </template>
      </el-table-column>>
      <el-table-column
        prop="images"
        label="镜像"
      >
        <template slot-scope="{row}">
          <span v-for="(tag,i) in row.images "><el-tag>{{ tag }}</el-tag></span>
        </template>
      </el-table-column>
      <el-table-column
        prop="node_name"
        label="调度节点"
      />
      <el-table-column
        prop="create_time"
        label="创建时间"
      />
      <el-table-column
        prop="restart_count"
        label="RESTARTS"
      />
      <el-table-column
        prop=""
        label="获取json"
      >
        <template slot-scope="scope">
          <el-button size="mini" type="success" @click="toPodjson(scope.row)">Get json</el-button>
        </template>
      </el-table-column>>
    </el-table>
  </div>
</template>

<script>
import { store } from '@/store'

const token = store.getters['user/token']
const user = store.getters['user/userInfo']
const path = process.env.VUE_APP_BASE_API
let timer = null

// import { useRoute } from 'vue-router'
// const route = useRoute()
export default {
  name: 'Detail',
  data() {
    return {
      tableData: [{}],
      pods: [{}],
      name: '',
      backData: {},
      my: [],
      a: 'skhfk',
      ns: ''
    }
  },
  computed: {},
  created() {
    console.log('启动')
    this.setName()
    this.getDeployDetail()
    this.setTimer()
  },
  mounted() {

  },
  methods: {
    toPodjson(row) {
      this.$router.push({
        name: 'podjson',
        params: {
        },
        query: {
          ns: row.ns,
          pod: row.name
        }
      })
    },
    deleteDeploy(row) {
      console.log(row.ns)
      return this.axios.delete(path + '/kubernetes/deploy/delete?' + 'ns=' + this.pods[0].ns + '&deploy=' + this.name, {
        headers: {
          'Content-Type': 'application/json',
          'x-token': token,
          'x-user-id': user.ID
        }
      })
        .then(rsp => {
          this.backData = rsp.data
          this.tableData = [{
            name: '',
            images_name: [
              '',
              ''
            ],
            replicas: 0,
            pods: [{ name: 'shfk' }, { name: 'shfk' }]
          }]
          this.tableData[0] = rsp.data
          this.pods = rsp.data.pods
          console.log(this.tableData)
        })
    },
    setName() {
      this.name = this.$route.query.name
    },
    handleStr: function(v, i) {
      if (i === 0) {
        return v
      } else if (i === 1) {
        return '+其他'
      } else {
        return ''
      }
    },
    getDeployDetail() {
      return this.axios.get(path + '/kubernetes/deploy/' + this.name, {
        headers: {
          'Content-Type': 'application/json',
          'x-token': token,
          'x-user-id': user.ID
        }
      })
        .then(rsp => {
          this.backData = rsp.data
          this.tableData = [{
            name: '',
            images_name: [
              '',
              ''
            ],
            replicas: 0,
            pods: [{ name: 'shfk' }, { name: 'shfk' }]
          }]
          this.tableData[0] = rsp.data
          this.pods = rsp.data.pods
          console.log(this.tableData)
        })
    },
    setTimer() {
      // eslint-disable-next-line no-undef
      this.axios.get(path + '/kubernetes/deploy/' + this.name, {
        headers: {
          'Content-Type': 'application/json',
          'x-token': token,
          'x-user-id': user.ID
        }
      }).then(rsp => {
        this.backData = rsp.data
        this.tableData = [{
          name: '',
          images_name: [
            '',
            ''
          ],
          replicas: 0,
          pods: [{ name: 'shfk' }, { name: 'shfk' }]
        }]
        this.tableData[0] = rsp.data
        this.pods = rsp.data.pods
        console.log(this.tableData)
        // eslint-disable-next-line no-constant-condition
        if (this.pods.length !== this.tableData.replicas) { // 根据返回状态判断
          timer = setTimeout(() => {
            this.setTimer()
          }, 2000)// 2秒查一下
        } else {
          clearTimeout(timer)// 清理定时任务
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
    toMy(v, i, sel) {
      if (i === sel && i !== 2) {
        return v + ' / '
      } else if (i === sel && i === 2) {
        return v + '      '
      } else {
        return ''
      }
    },
    incReplica(ns, deploy, dec) {
      return this.axios({
        method: 'post',
        url: path + '/kubernetes/deploy/replicas',
        headers: {
          'Content-Type': 'application/json',
          'x-token': token,
          'x-user-id': user.ID
        },
        data: {
          ns: ns,
          deploy: deploy,
          dec: dec
        }
      })
        .then(rsp => {
          console.log(rsp)
        })
    }

  }

}
</script>

<style scoped>
.bg {
  background-color: white;
}

.ft {
  color: black;
}

.red {
  color: red;
}

.green {
  color: green;
}

.black {
  color: black;
}

.bigWord {
  font-size: 30px;
}
</style>
