const sql = {
    state: {
      orders: {
        my_order: localStorage.getItem("yasql.orders.my_order") || ""
      },
      sqlquery: {
        query_hash: localStorage.getItem("yasql.sqlquery.query_hash") || "",
        user_input: localStorage.getItem("yasql.sqlquery.user_input") || "SELECT ",
        user_theme: localStorage.getItem("yasql.sqlquery.user_theme") || "default"
      }
    },
    mutations: {
      STORE_MY_ORDER: (state, value) => {
        state.orders.my_order = value;
        localStorage.setItem("yasql.orders.my_order", value);
      },
      STORE_QUERY_HASH: (state, value) => {
        state.sqlquery.query_hash = value;
        localStorage.setItem("yasql.sqlquery.query_hash", value);
      },
      STORE_USER_INPUT: (state, value) => {
        state.sqlquery.user_input = value;
        localStorage.setItem("yasql.sqlquery.user_input", value);
      },
      STORE_USER_THEME: (state, value) => {
        state.sqlquery.user_theme = value;
        localStorage.setItem("yasql.sqlquery.user_theme", value);
      }
    },
    actions: {
      storeMyOrder({ commit }, playload) {
        commit("STORE_MY_ORDER", playload);
      },
      storeQueryHash({ commit }, playload) {
        commit("STORE_QUERY_HASH", playload);
      },
      storeUserIuput({ commit }, playload) {
        commit("STORE_USER_INPUT", playload);
      },
      storeUserTheme({ commit }, playload) {
        commit("STORE_USER_THEME", playload);
      }
    }
  };
  
  export default sql;
  