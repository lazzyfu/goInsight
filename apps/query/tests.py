aa = "business_tasks&I_BT_ID&int(10) unsigned,business_tasks&I_TYPE&mediumint(10) unsigned,business_tasks&I_OWNER_ID&int(10) unsigned,business_tasks&I_AMOUNT&decimal(11,2),business_tasks&I_PROJECT_ID&int(10) unsigned,business_tasks&CH_PROJECT_NAME&varchar(100),business_tasks&D_FINISHED_AT&datetime,business_tasks&D_CREATED_AT&datetime,business_tasks&D_UPDATED_AT&datetime"

for i in aa.split(','):
    print(i.strip())