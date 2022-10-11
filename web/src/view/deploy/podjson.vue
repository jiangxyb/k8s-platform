<template>
  <div>
    <pre>{{this.podjson}}</pre>
  </div>
</template>

<script>
import { store } from '@/store'

const token = store.getters['user/token']
const user = store.getters['user/userInfo']
const path = process.env.VUE_APP_BASE_API

// import { useRoute } from 'vue-router'
// const route = useRoute()
export default {
  name: 'Podjson',
  data() {
    return {
      ns: '',
      pod: '',
      podjson: ''
    }
  },
  computed: {},
  created() {
    console.log('启动')
    this.setNsAndPod()
    this.getPodJson()
  },
  mounted() {

  },
  methods: {
    setNsAndPod() {
      this.ns = this.$route.query.ns
      this.pod = this.$route.query.pod
    },
    getPodJson() {
      return this.axios.get(path + '/kubernetes/pod/json?' + 'ns=' + this.ns + '&pod=' + this.pod, {
        headers: {
          'Content-Type': 'application/json',
          'x-token': token,
          'x-user-id': user.ID
        }
      })
        .then(rsp => {
          this.podjson = JSON.stringify(rsp, null, 4)
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
