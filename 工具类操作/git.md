# git操作

```bash
git remote remove  origin

echo "# k8s-study" >> README.md
git init
git commit -m "first commit"
git remote add origin git@github.com:luoguoling/k8s-study.git
git branch -M main
git push -u origin main

查看分支
git branch -a
在本地仓库删除文件
git rm 我的文件

在本地仓库删除文件夹
git rm -r 我的文件夹/

排除某个目录
 git submodule add <url> k8s-study  增加

git rm --cached k8s-study
放弃本地修改
git checkout .
放弃某个文件修改
git checkout -- filepathname
切换分支
git checkout main


回滚
git log -x 查看最新的几个历史版本信息
git reset --hard xxx  回退到具体版本号
git reset --hard HEAD^ 回退到上个版本

git reset --hard HEAD~3 回退到前3次提交之前，以此类推，回退到n次提交之前

git reset --hard commit_id 退到/进到，指定commit的哈希码（这次提交之前或之后的提交都会回滚）

回滚之后强制推到远程
git push origin HEAD --force

合并分支
1，将开发分支代码合入到master中
git checkout dev           #切换到dev开发分支
git pull
git checkout master
git merge dev              #合并dev分支到master上
git push origin master     #将代码推到master上

2.将master的代码同步更新 到开发分支中
git checkout master
git pull
git checkout dev
git merge master
git pull origin dev
```

> https://git-scm.com/book/zh/v2/Git-%E5%88%86%E6%94%AF-%E5%88%86%E6%94%AF%E7%9A%84%E6%96%B0%E5%BB%BA%E4%B8%8E%E5%90%88%E5%B9%B6#