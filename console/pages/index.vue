<template>
  <div class="console" :style="flex" id="particles">
    <div class="card ft-center" ref="content">
      <h2 class="card__title">{{ $t('popularBlog', $store.state.locale) }}</h2>
      <ul class="list ">
        <li class="fn-flex" v-for="item in list">
          <a class="fn-flex-1" :href="item.url" target="_blank">{{ item.title }}</a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
  import 'particles.js'
  import Vue from 'vue'
  import { initParticlesJS } from '~/plugins/utils'

  export default {
    head () {
      return {
        title: this.$t('welcome', this.$store.state.locale)
      }
    },
    data () {
      return {
        flex: '',
        list: []
      }
    },
    async mounted () {
      const responseTopData = await this.axios.get('/blogs/top')
      if (responseTopData) {
        this.$set(this, 'list', responseTopData)
        Vue.nextTick(() => {
          if (this.$refs.content.scrollHeight > this.$refs.content.clientHeight) {
            this.$set(this, 'flex', 'flex:none')
          }
        })
      }
      initParticlesJS('particles')
    }
  }
</script>

<style lang="sass">
  .card__title
    display: block
</style>
