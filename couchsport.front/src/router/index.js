import Vue from 'vue'
import Router from 'vue-router'

import AuthPages from 'pages/auth/AuthPages'
import About from 'pages/About'
import Home from 'pages/Home'
import Explore from 'pages/Explore'
import Profile from 'pages/profile/Profile'
import PageDetails from 'pages/PageDetails'

import Informations from 'pages/profile/Informations'
import Activities from 'pages/profile/Activities'
import Pages from 'pages/profile/Pages'
import Conversations from 'pages/profile/Conversations'

import { defaultLocale, i18n } from '../trans'

Vue.use(Router)

var router

export default function (store) {
  if (router) return router

  const requiresToLogin = (to, from, next) => {
    if (!store.getters.isAuthenticated) {
      next()
      return
    }
    next({ name: 'informations', params: { locale: i18n.locale } })
  }

  const requiresPrivileges = (to, from, next) => {
    if (store.getters.isAuthenticated) {
      next()
      return
    }
    next({ name: 'login', params: { locale: i18n.locale } })
  }

  const originalPush = Router.prototype.push

  Router.prototype.push = function push(location) {
    return originalPush.call(this, location).catch(() => console.log)
  }

  const routes = {
    mode: 'history',
    routes: [
      {
        path: '/',
        redirect: `/${store.getters.getLocale || defaultLocale}`
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
            beforeEnter: requiresToLogin
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
            beforeEnter: requiresPrivileges,
            children: [
              {
                path: '',
                name: 'informations',
                component: Informations
              },
              {
                path: 'activities',
                name: 'activities',
                component: Activities
              },
              {
                path: 'pages',
                name: 'pages',
                component: Pages
              },
              {
                path: 'conversations',
                name: 'conversations',
                component: Conversations
              }
            ]
          }
        ]
      }
    ]
  }

  router = new Router(routes)

  router.beforeEach((to, from, next) => {
    let language = to.params.locale
    if (['fr', 'en'].indexOf(language) < 0) {
      language = store.getters.getLocale || defaultLocale
    }

    i18n.locale = language
    to.params.locale = language
    next()
  })

  return router
}
