package util

const (
	//table system
	TABLE_AUTO_SYS_IMAGE                   = "auto_sys_image"
	TABLE_AUTO_SYS_FLOW                    = "auto_sys_flow"
	TABLE_AUTO_SYS_PARAMETER               = "auto_sys_parameter"
	TABLE_AUTO_SYS_JOBSPOOL                = "auto_sys_jobspool"
	TABLE_AUTO_SYS_SLAVE_HEART             = "auto_sys_slave_heart"
	TABLE_AUTO_SYS_WORKER_HEART            = "auto_sys_worker_heart"
	TABLE_AUTO_SYS_MASTER_HEART            = "auto_sys_master_heart"
	TABLE_AUTO_SYS_MASTER_ROUTINE_HEART    = "auto_sys_master_routine_heart"
	TABLE_AUTO_SYS_FLOW_MASTER             = "auto_sys_flow_master"
	TABLE_AUTO_SYS_JOB_RUNNING_HEART       = "auto_sys_job_running_heart"
	TABLE_AUTO_SYS_SLAVE_JOB_RUNNING_HEART = "auto_sys_slave_job_running_heart"
	//
	TABLE_AUTO_FLOW_PARAMETER = "auto_flow_parameter"
	//table flow job
	TABLE_AUTO_JOB            = "auto_job"
	TABLE_AUTO_JOB_DEPENDENCY = "auto_job_dependency"
	TABLE_AUTO_JOB_TIMEWINDOW = "auto_job_timewindow"
	TABLE_AUTO_JOB_CMD        = "auto_job_cmd"
	TABLE_AUTO_JOB_SLAVE      = "auto_job_slave"
	TABLE_AUTO_JOB_PARAMETER  = "auto_job_parameter"
	TABLE_AUTO_JOB_STREAM     = "auto_job_stream"
	TABLE_AUTO_JOB_LOG        = "auto_job_log"
	//file
	FILE_AUTO_SYS_DBSTORE = "auto-dbstore.db"
	//status
	STATUS_AUTO_READY   = "Ready"
	STATUS_AUTO_RUNNING = "Running"
	STATUS_AUTO_PENDING = "Pending"
	STATUS_AUTO_STOP    = "Stop"
	STATUS_AUTO_SUBMIT  = "Submit"
	STATUS_AUTO_GO      = "Go"
	STATUS_AUTO_FAIL    = "Fail"
	STATUS_AUTO_SUCC    = "Succ"
	//
	CONST_ENABLE   = "1"
	CONST_DISABLE  = "0"
	CONST_FLOW_CTS = "AUTO_CREATE_TIMESTAMP"
	CONST_FLOW_RCT = "AUTO_RUN_CONTEXT"
	CONST_SYS      = "AUTO_SYS"
	CONST_JOB      = "AUTO_JOB"
	CONST_FLOW     = "AUTO_FLOW"
)
