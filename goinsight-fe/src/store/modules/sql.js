const sql = {
    state: {
      orders: {
        my_order: localStorage.getItem("orders.my_order") || ""
      },
    },
    mutations: {
      STORE_MY_ORDER: (state, value) => {
        state.orders.my_order = value;
        localStorage.setItem("orders.my_order", value);
      },
    },
    actions: {
      storeMyOrder({ commit }, playload) {
        commit("STORE_MY_ORDER", playload);
      },
    }
  };
  
  export default sql;
  