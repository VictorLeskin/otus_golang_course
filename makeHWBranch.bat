set HW=hw02_unpack_string
set HW=hw03_frequency_analysis
set HW=hw06_pipeline_execution
set HW=hw07_file_copying
set HW=hw08_envdir_tool
set HW=hw10_program_optimization

git switch master
pause

git checkout -b %HW%
# create join the remove and local branchs
git branch --set-upstream-to=origin/%HW% %HW%