set HW=hw02_unpack_string
set HW=hw03_frequency_analysis
set HW=hw06_pipeline_execution
set HW=hw09_struct_validator

git switch master
pause

git checkout -b %HW%
cd %HW%
sed 's/fixme_my_friend/VictorLeskin\/otus_golang_course/' -i go.mod 
rm -rf .sync

git commit -am "Initial commit: fix web path of project. Remove .sync file"

# create branch on github and push the commit
git push origin %HW%:%HW%

# create join the remove and local branchs
git branch --set-upstream-to=origin/%HW% %HW%