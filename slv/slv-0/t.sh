
t="
CREATE TABLE dw_sdk.ext_thirdapp_activedevice_d_incr(
  product_id string, 
  dhid string, 
  vercode string, 
  chanid string, 
  province string, 
  open_cnt string, 
  fcv string, 
  fchannel string, 
  fprovince string, 
  fpt string, 
  aid string, 
  old_dhid string, 
  longi string, 
  lati string)
PARTITIONED BY ( 
  pt string)
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.ql.io.orc.OrcSerde' 
STORED AS INPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.orc.OrcInputFormat' 
OUTPUTFORMAT 
  'org.apache.hadoop.hive.ql.io.orc.OrcOutputFormat'
LOCATION
  'hdfs://ns1/user/hive/warehouse/sdk/dw_sdk.db/ext_thirdapp_activedevice_d_incr'
TBLPROPERTIES (
  'last_modified_by'='sdk', 
  'last_modified_time'='1558518192', 
  'transient_lastDdlTime'='1558518192');
"
echo "$t"
echo "$t"
sleep 30
echo "ok"
echo "$t"
