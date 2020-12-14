const route = {
    name: "view.sqlquery",
    path: "/sqlquery",
    component: () => import("./index.vue"),
    meta: { title: "DMS", keepAlive: true, icon: "search" },
    children: [
      {
        name: "view.sqlquery.query",
        path: "/sqlquery/query",
        component: () => import("./query/index.vue"),
        meta: { title: "DB查询" }
      },
    ]
  };
  
  export default route;
  