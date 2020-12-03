import Vue from 'vue'
import Router from 'vue-router'

import Login from 'components/auth/Login'
import SignUp from 'components/auth/SignUp'
import About from 'components/About'
import Home from 'components/Home'
import Explore from 'components/explore/Explore'
import Profile from 'components/profile/Profile'
import PageDetails from 'components/page/PageDetails'

import { defaultLocale, i18n } from '../trans'
import store from '../store'

Vue.use(Router)

const ifNotAuthenticated = (to, from, next) => {
  if (!store.getters.isAuthenticated) {
    next()
    return
  }
  next({ name: 'profile', params: { locale: i18n.locale } })
}

const ifAuthenticated = (to, from, next) => {
  if (store.getters.isAuthenticated) {
    next()
    return
  }
  next({ name: 'login', params: { locale: i18n.locale } })
}

const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch((err) => err)
}

let router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      redirect: `/${defaultLocale}`
    },
    {
      path: '/:locale',
      component: {
        template: '<router-view />'
      },
      children: [
        {
          path: '/',
          name: 'home',
          component: Home,
          props: true
        },
        {
          path: 'explore',
          name: 'explore',
          component: Explore
        },
        {
          path: 'signup',
          name: 'signup',
          component: SignUp
        },
        {
          path: 'login',
          name: 'login',
          component: Login,
          props: true,
          beforeEnter: ifNotAuthenticated
        },
        {
          path: 'pages/:page_name',
          name: 'page-details',
          component: PageDetails
        },
        {
          path: 'about',
          name: 'about',
          component: About
        },
        {
          path: 'profile',
          name: 'profile',
          component: Profile,
          beforeEnter: ifAuthenticated
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  let language = to.params.locale
  if (!language) {
    language = defaultLocale
  }

  i18n.locale = language
  next()
})

export default router
