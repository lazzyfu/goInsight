# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.urls import path

from sqlorders import views

urlpatterns = [
    # SQL工单
    path('envs', views.GetDBEnvironment.as_view(), name='v1.sqlorders.db-environment'),
    path('schemas', views.GetDbSchemas.as_view(), name='v1.sqlorders.db-schemas'),
    path('incep/syntaxcheck', views.IncepSyntaxCheckView.as_view(), name='v1.sqlorders.incep.syntaxcheck'),
    path('commit', views.SqlOrdersCommit.as_view(), name='v1.sqlorders.commit'),
    path('list', views.SqlOrdersList.as_view(), name='v1.sqlorders.list'),
    path('detail/<str:order_id>', views.SqlOrdersDetail.as_view(), name='v1.sqlorders.detail'),
    path('op/approve/<int:pk>', views.OpSqlOrderView.as_view({"put": "approve"}), name='v1.sqlorders.approve'),
    path('op/feedback/<int:pk>', views.OpSqlOrderView.as_view({"put": "feedback"}), name='v1.sqlorders.feedback'),
    path('op/close/<int:pk>', views.OpSqlOrderView.as_view({"put": "close"}), name='v1.sqlorders.close'),
    path('op/review/<int:pk>', views.OpSqlOrderView.as_view({"put": "review"}), name='v1.sqlorders.review'),
    # 生成工单任务
    path('tasks/generate', views.GenerateTasksView.as_view(), name='v1.sqlorders.generate-tasks'),
    path('tasks/get/<str:order_id>', views.GetTaskIdView.as_view(), name='v1.sqlorders.get-task-id'),
    path('tasks/list/<str:task_id>', views.GetTasksListView.as_view(), name='v1.sqlorders.get-tasks-list'),
    path('tasks/preview/<str:task_id>', views.GetTasksPreviewView.as_view(),
         name='v1.sqlorders.get-tasks-preview'),
    # 执行任务
    path('tasks/execute/single', views.ExecuteSingleTaskView.as_view(), name='v1.sqlorders.execute-single-task'),
    path('tasks/execute/multi', views.ExecuteMultiTasksView.as_view(), name='v1.sqlorders.execute-multi-tasks'),
    path('tasks/throttle', views.ThrottleTaskView.as_view(), name='v1.sqlorders.throttle-task'),
    path('tasks/result/<int:id>', views.GetTasksResultView.as_view(), name='v1.sqlorders.get-tasks-result'),
    # Hook
    path('hook', views.HookSqlOrdersView.as_view(), name='v1.sqlorders.hook-sqlorders'),
    # download export files
    path('export/download/<str:base64_filename>', views.DownloadExportFilesView.as_view(),
         name='v1.sqlorders.download-export-files'),
    # 上线版本
    path('versions/get', views.ReleaseVersionsGet.as_view(), name='v1.sqlorders.versions.get'),
    path('versions/list', views.ReleaseVersionsList.as_view(), name='v1.sqlorders.versions.list'),
    path('versions/create', views.ReleaseVersionsCreate.as_view(),
         name='v1.sqlorders.versions.create'),
    path('versions/update/<int:key>', views.ReleaseVersionsUpdate.as_view(),
         name='v1.sqlorders.versions.update'),
    path('versions/delete/<int:id>', views.ReleaseVersionsDelete.as_view(),
         name='v1.sqlorders.versions.delete'),
    path('versions/view/<str:version>', views.ReleaseVersionsView.as_view(),
         name='v1.sqlorders.versions.view'),
]
