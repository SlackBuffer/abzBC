<!-- https://bitcoinbook.cs.princeton.edu/ -->
- **Credit** and **cash** are fundamental ideas, to the point we can sort the multitude of electronic payment methods into 2 piles

- Credit
    - Intermediary architecture (e.g., SET architecture)
        - CyberCash implemented SET. In addition to credit card payment processing, they had a digital cash product called CyberCoin intended for **small payments** such as paying a few cents to read an online newspaper article. It was affected by the Y2K bug - it caused their payment processing software to double-bill some customers
        - The fundamental problem with SET has to do with ***certificates***
- A certificate is a way to securely ***associate*** a cryptographic identity, that is, a public key, with a real-life identity
    - > It’s what a website needs to obtain, from companies like Verisign that are called certification authorities, in order to show up as secure in your browser (typically indicated by a lock icon)
    - Putting security before usability, CyberCash and SET decided that not only would processors and merchants in their system have to get certificates, all users would have to get one as well
    - Getting a certificate is about as pleasant as doing your taxes, so the system was a disaster. Over the decades, mainstream users have said a firm and collective ‘no’ to any system that requires end-user certificates
    - In Bitcoin, public keys themselves are the identities by which users are known

- Cash
    - A cash system needs to be "bootstrapped" (要有初始的货币分配)
        1. Avoids the possibility of a buyer defaulting on his debt
        2. Better anonymity
        3. Enables offline transactions where there's no need to phone home to a third party in order to get the transaction approved
    - Bitcoin is not anonymous to the same level as cash is
        - You don’t need to use your real identity to pay in Bitcoin, but it’s possible that your transactions can be tied together based on the public ledger of transactions with clever algorithms, and then further linked to your identity if you’re not careful
    - Bitcoin doesn’t work in a fully offline way either
        - It doesn’t require a central server, instead relying on a peer-to-peer network which is resilient in the way that the Internet itself is
        - Tricks like “green addresses” and micropayments which allow us to do offline payments in certain situations or under certain assumptions
    - Double spending solution
        1. 实体票据
            1. 权威签名（为兑换货币背书
            2. 序列号（防伪造后的重复消费企图）
                - 票据兑换信息记录在票据发行方
                - 接收票据前需跟实体票据发行方确认该票据是否已被兑换
                    - 收到票据的一方得到确认票据未被兑换，随后接收票据，同时发行方将票据的状态改为已兑换
                    - 接收方若想使用该票据，需到发行方处兑换等额的序列号不同的新票据，才能使用
            - 此种方案**失去了货币的匿名性**特点，除了记录序列号，发行方在发行票据和他人前来兑换时均可记录用户的身份
        2. 票据方案改进一
            1. 发行方发行票据给用户时，用户为该票据挑选序列号且不告知发行方，发行方盲签该票据（blind signature)
            2. 每个用户都会自发的去选择的尽量长的随机序列号以使得自己的序列号唯一
            - 缺点是需要每个人都信任的中心化权威机构提供的中心化服务器来支撑交易
        3. 票据方案改进二
            - 重点放在**检测**重复消费，而不是去阻止；检测到之后追回损失、惩戒作恶者；若数字货币被广泛接纳，法律系统便会将重复消费视为犯罪行为
            - 密码学原理
                1. 发行给用户的数字货币加密了用户的身份后生成密文，只有用户能解
                2. 每次消费时，收银员要求用户解密密文的随机子串，并记录解密得到的结果，该结果不足以确定用户的身份
                3. 若用户存在重复消费消费行为，最终两位收银员都去银行兑换该用户支付的数字货币，此时银行就能通过整合 2 段解密后的信息得到用户的身份
                4. 第一位收银员通过重复消费用户的数字货币以陷害之的企图行不通：第二位收银员会要求陷害者解密该数字货币密文的随机子串，而该子串与第一位收银员的子串相同的概率微乎其微，所以第二位收银员无法成功解密子串，因而无法完成交
- 数字货币找零问题
    - Merkle trees
- Solutions to computational puzzles could be digital objects that have some value
    - > Email spam
- Block chain
    - The initial proposal was a method for secure timestamping of digital documents
        - 时间戳服务器接收到文档后，对文档、当前时间和指向前一个文档的链接（hash pointer）这三个内容签名，并为该签名颁发证书；三个内容任一变化后，链接自动失效
    - 优化：将文档打成 block 后再连接；block 内的文档是树形连接而不是线性连接
    - Hashcash-esque 协议用于调控出块速度