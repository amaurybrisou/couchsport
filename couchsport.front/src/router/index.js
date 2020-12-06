import Vue from 'vue'
import Router from 'vue-router'

import AuthPages from 'pages/auth/AuthPages'
import About from 'pages/About'
import Home from 'pages/Home'
import Explore from 'pages/Explore'
import Profile from 'pages/profile/Profile'
import PageDetails from 'pages/PageDetails'

import { defaultLocale, i18n } from '../trans'

Vue.use(Router)

export default function (store) {
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

  const routes = {
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
            component: AuthPages
          },
          {
            path: 'login',
            name: 'login',
            component: AuthPages,
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
  }

  let router = new Router(routes)

  router.beforeEach((to, from, next) => {
    let language = to.params.locale
    if (!language) {
      language = defaultLocale
    }

    i18n.locale = language
    next()
  })

  return router
}
