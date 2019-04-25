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