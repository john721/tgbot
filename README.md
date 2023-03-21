# 幫選晚餐機器人

家裡小孩上了高中<br />
在學校下課之後到社團練習或去補習班之前<br />
有一段空檔可以吃晚餐<br />
總是花了一半時間猶豫要吃什麼<br />
"希望有程式幫忙選晚餐, 跳過昨天和前天吃過的"<br />
這個 Telegram Bot 是醬子來的..<br />

# 功能
- 建立個人餐點列表
- 如果菜單 > 3種, 會選擇昨天和前天沒吃過的
- 廢到讓你想笑的陽春功能

![看起來像醬](screenshot.png)

# 怎麼弄?
1. 打開 Telegram, 加入 @BotFather, 或是點一下這裡 https://t.me/BotFather
2. 和 BotFather 對話, 用指令 /newbot 可以新建一個機器人, 它會問這隻機器人基本設定.
3. 完成之後, BotFather 會給你一個 Token, 長得像是醬子: <br />
   1234567890:ABCDEF_abcdefghijklmnopqrstu-VWXYZA
4. 把取得的 token 設為環境變數:<br />
   export TELEGRAM_APITOKEN="1234567890:ABCDEF_abcdefghijklmnopqrstu-VWXYZA"
5. ./run.sh