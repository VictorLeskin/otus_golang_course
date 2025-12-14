set HW=hw02_unpack_string
set HW=hw03_frequency_analysis

git switch master

git checkout -b %HW%
cd %HW%
sed 's/fixme_my_friend/VictorLeskin\/otus_golang_course/' -i go.mod 
rm -rf .sync

git commit -am "Initial commit: fix web path of project. Remove .sync file"

# create branch on github and push the commit
git push origin %HW%:%HW%

# create join the remove and local branchs
git branch --set-upstream-to=origin/%HW% %HW%