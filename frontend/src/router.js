import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
	routes: [
		{
			path: "/",
			redirect: "/books",
		},
		{
			path: '/books',
			name: 'list-book',
			component: () => import(/* webpackChunkName: "books" */ './views/ListBook.vue')
		},
		{
			path: '/add',
			name: 'add',
			component: () => import(/* webpackChunkName: "add" */ './views/EditBook.vue')
		},
		{
			path: '/edit',
			name: 'edit',
			component: () => import(/* webpackChunkName: "edit" */ './views/EditBook.vue'),
			props: true,
		},
		{
			path: '/books/:id',
			name: 'detail-book',
			component: () => import(/* webpackChunkName: "detail" */ './views/DetailBook.vue')
		},
		{
			path: '/histories',
			name: 'list-record',
			component: () => import(/* webpackChunkName: "histories" */ './views/ListRecord.vue')
		},
		{
			path: '/admin',
			name: 'admin',
			component: () => import(/* webpackChunkName: "admin" */ './views/Admin.vue'),
		},
	]
})
