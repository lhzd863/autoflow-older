Sub job()
'功能:生成作业配置'
'DATE:2020-06-12  '
'#####################常量定义区#####################'

 FirstLine = 2               '配置信息起始行'
 Firstcol = 1                '配置信息起始列'
 
 '--------------------------------------------------'
 

'#####################默认值#####################'
 Status = "Ready"
 JobType = "D"
 Enable = 1
 RetryCnt = 0
 RetryCnt = 0
 Alert = "N"
 TimeTrigger = "N"
 Frequency = "N"
 CheckBatStatus = "N"
 Priority = 1
 '####################################################'

If ActiveSheet.Name <> "job" Then
  MsgBox "此宏在“job生成页”执行", vbInformation, "提示"
  Exit Sub
End If

Open ThisWorkbook.Path & "\" & "job.yaml" For Output As #1                       'script context'

arr = Range("A1:N" & Range("A65535").End(xlUp).Row)

introw = UBound(arr)

MsgBox introw, vbInformation, "提示"

Print #1, "apiversion: v1"
Print #1, "name: job"
Print #1, "labels:"

For i = FirstLine To UBound(arr)
   Print #1, "  #" & arr(i, Firstcol + 3)
   Print #1, "  " & UCase(Trim(arr(i, Firstcol + 1))) & "." & UCase(Trim(arr(i, Firstcol + 2))) & ":"
   Print #1, "    sys: """ & UCase(Trim(arr(i, Firstcol + 1))) & """"
   Print #1, "    job: """ & UCase(Trim(arr(i, Firstcol + 2))) & """"
   Print #1, "    runcontext: """""
   Print #1, "    status: """ & Status & """"
   Print #1, "    starttime: """""
   Print #1, "    endtime: """""
   Print #1, "    jobtype: """ & JobType & """"
   Print #1, "    description: """ & Trim(arr(i, Firstcol + 3)) & """"
   Print #1, "    enable: """ & Enable & """"
   Print #1, "    sserver: """""
   Print #1, "    sip: """""
   Print #1, "    sport: """""
   Print #1, "    mserver: """""
   Print #1, "    mip: """""
   Print #1, "    mport: """""
   Print #1, "    routineid: """""
   Print #1, "    timewindow: """""
   Print #1, "    retrycnt: """ & RetryCnt & """"
   Print #1, "    alert: """ & Alert & """"
   Print #1, "    timetrigger: """ & TimeTrigger & """"
   Print #1, "    frequency: """ & Frequency & """"
   Print #1, "    checkbatstatus: """ & CheckBatStatus & """"
   Print #1, "    priority: """ & Priority & """"
   Print #1, "    runningcmd: """""
   Print #1, "    isreplace: """ & arr(i, Firstcol + 0) & """"
   Print #1, ""
Next i
 
Close #1
End Sub

