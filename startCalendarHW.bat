set HW=hw02_unpack_string
set HW=hw03_frequency_analysis
set HW=hw06_pipeline_execution
set HW=hw07_file_copying
set HW=hw08_envdir_tool
set HW=hw09_struct_validator
set HW=hw10_program_optimization
set HW=hw11_telnet_client
set HW=hw12_calendar
set HW=hw13_calendar

git switch hw12_calendar
pause

git checkout -b %HW%

# create join the remove and local branchs
git push -u origin %HW%
