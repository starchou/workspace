dns
==========


动态根据宽带public ip更新dnspod登记的域名
按照 https://gist.github.com/833369 逻辑重新用Go实现了，用更少的内存开销在Raspberry Pi上跑。

替换上你的Email，密码，域名ID，记录ID等参数，就可以运行了。 会在后台一直运行，每隔30秒检查一遍IP，如果修改了就更新IP。

获得domain_id可以用：

curl curl -k https://dnsapi.cn/Domain.List -d "login_email=xxx&login_password=xxx" 

获得record_id：

curl -k https://dnsapi.cn/Record.List -d "login_email=xxx&login_password=xxx&domain_id=xxx"

iplocation
==========

根据python版本 改编而来
详细原理：http://lumaqq.linuxsir.org/article/qqwry_format_detail.html


oauth
==========

微博，QQ，豆瓣，人人 社會化登陸API接口