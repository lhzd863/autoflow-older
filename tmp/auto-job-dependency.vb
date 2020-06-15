Sub depenedency()
'功能:生成作业配置'
'DATE:2020-06-12  '
'#####################常量定义区#####################'

 FirstLine = 2               '配置信息起始行'
 Firstcol = 1                '配置信息起始列'
 
 '--------------------------------------------------'
 

'#####################默认值#####################'
 Enable = 1
 '####################################################'

If ActiveSheet.Name <> "dependency" Then
  MsgBox "此宏在“job生成页”执行", vbInformation, "提示"
  Exit Sub
End If

Open ThisWorkbook.Path & "\" & "dependency.yaml" For Output As #1                       'script context'

arr = Range("A1:N" & Range("A65535").End(xlUp).Row)

introw = UBound(arr)

MsgBox introw, vbInformation, "提示"

Print #1, "apiversion: v1"
Print #1, "name: dependencyjob"
Print #1, "labels:"

For i = FirstLine To UBound(arr)
   Print #1, "  #" & arr(i, Firstcol + 3)
   Print #1, "  " & UCase(Trim(arr(i, Firstcol + 1))) & "." & UCase(Trim(arr(i, Firstcol + 2))) & ":"
   Print #1, "    sys: """ & UCase(Trim(arr(i, Firstcol + 1))) & """"
   Print #1, "    job: """ & UCase(Trim(arr(i, Firstcol + 2))) & """"
   Print #1, "    dependencysys: """ & UCase(Trim(arr(i, Firstcol + 3))) & """"
   Print #1, "    dependencyjob: """ & UCase(Trim(arr(i, Firstcol + 4))) & """"
   Print #1, "    enable: """ & Enable & """"
   Print #1, "    isreplace: """ & arr(i, Firstcol + 0) & """"
   Print #1, ""
Next i
 
Close #1
End Sub


