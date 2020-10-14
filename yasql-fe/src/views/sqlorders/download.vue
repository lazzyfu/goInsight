<template>
    <a-card title="下载导出文件">
        <exception-page :type="type" />
    </a-card>
</template>

<script>
import { downloadExportFiles } from "@/api/sql";
import ExceptionPage from "./exception";

export default {
    components: {
        ExceptionPage,
    },
    data() {
        return {
            file_name: "",
            type: 200,
        };
    },
    methods: {
        download() {
            const base64_filename = this.$route.params.base64_filename;
            this.file_name = atob(base64_filename);
            downloadExportFiles(base64_filename)
                .then((response) => {
                    let blob = new Blob([response], { type: "application/zip" });
                    let url = URL.createObjectURL(blob);
                    let link = document.createElement("a");
                    link.href = url;
                    link.download = this.file_name;
                    link.click();
                    URL.revokeObjectURL(url);
                })
                .catch((error) => {
                    console.log('error: ', error);
                    // 此处的403不工作，是被全局定义的axios.interceptors.response.use拦截了
                    const errors = [403, 404];
                    if (errors.includes(error.response.status)) {
                        this.type = error.response.status;
                    }
                });
        },
    },
    mounted() {
        this.download();
    },
};
</script>

<style>
</style>