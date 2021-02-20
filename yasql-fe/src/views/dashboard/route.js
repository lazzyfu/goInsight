const route = {
    name: "view.dashboard",
    path: "/dashboard",
    component: () => import("./index.vue"),
    meta: { title: "仪表盘", keepAlive: true, icon: "search" }
  };
  
  export default route;
  