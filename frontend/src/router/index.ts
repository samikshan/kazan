import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "home",
    component: Home,
    meta: {
      title: 'Home',
      metaTags: [
        {
          name: 'description',
          content: 'Home page.'
        },
        {
          property: 'og:description',
          content: 'Home page.'
        }
      ]
    }
  },
  {
    path: '/signin',
    name: 'signin',
    component: () => import('@/views/Signin.vue'),
    meta: {
      title: 'Sign Into Your Kazan Account',
      metaTags: [
        {
          name: 'description',
          content: 'Sign In.'
        },
        {
          property: 'og:description',
          content: 'Sign In.'
        }
      ]
    }
  },
  {
    path: '/signup',
    name: 'signup',
    component: () => import('@/views/Signup.vue'),
    meta: {
      title: 'Sign Up With Kazan',
      metaTags: [
        {
          name: 'description',
          content: 'Sign Up.'
        },
        {
          property: 'og:description',
          content: 'Sign Up.'
        }
      ]
    }
  },
  {
    path: '/record',
    name: 'record',
    component: () => import('@/views/Recorder.vue'),
    meta: {
      title: 'Publish New Tracks',
      metaTags: [
        {
          name: 'description',
          content: 'Page to record and upload tracks.'
        },
        {
          property: 'og:description',
          content: 'Page to record and upload tracks.'
        }
      ]
    }
  },
  {
    path: '/jam',
    name: 'jam',
    component: () => import('@/views/Mixer.vue'),
    meta: {
      title: 'Jam On Track',
      metaTags: [
        {
          name: 'description',
          content: 'Page to jam on existing tracks.'
        },
        {
          property: 'og:description',
          content: 'Page to jam on existing tracks.'
        }
      ]
    }
  },
  {
    path: '/search/:text',
    name: 'search',
    component: () => import('@/views/Search.vue'),
    
  },
  
  {
    path: '/:id',
    name: 'profile',
    component: () => import('@/views/Profile.vue'),
  },
];

const router = new VueRouter({
  routes
});

export default router;
