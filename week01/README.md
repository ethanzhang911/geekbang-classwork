Q:我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？



A:应该Wrap这个error，抛给上层，但是因为是返回空行。针对这样的第三方库或者标准库，应该是用诸如pkg/erros等将error错误Wrap起来，抛给上层。





数据库如下

![image-20210725235133468](https://typorabyethancheung911.oss-cn-shanghai.aliyuncs.com/typora/image-20210725235133468.png)

如果按照id=1查找，应该返回空行。

