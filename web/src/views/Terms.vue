<template>
  <section class="view-terms">
    <h1 class="page-title">服务条款</h1>
    <div class="terms-card card">
      <p class="terms-updated">最后更新：2025年1月</p>

      <div class="terms-content">
        <h2>一、接受条款</h2>
        <p>
          欢迎使用衣搭库（以下简称「本网站」或「我们」）。使用本网站即表示您同意遵守本服务条款。若您不同意本条款，请勿使用本网站及相关服务。
        </p>

        <h2>二、服务说明</h2>
        <p>
          衣搭库为用户提供个人衣橱管理、穿搭推荐、收藏与分享等功能。我们保留在不另行通知的情况下修改、暂停或终止部分或全部服务的权利。
        </p>

        <h2>三、用户账号与责任</h2>
        <p>您应当：</p>
        <ul>
          <li>提供真实、准确、完整的注册信息；</li>
          <li>妥善保管账号与密码，对账号下的行为负责；</li>
          <li>不得将账号出借、转让或以其他方式供他人使用；</li>
          <li>不得利用本网站从事违法违规或侵犯他人权益的行为。</li>
        </ul>

        <h2>四、用户内容</h2>
        <p>
          您上传的衣物图片、搭配记录等（「用户内容」）仍归您所有。您授予我们为提供、改进服务所必需的、非独占的、可再许可的使用权。您承诺对用户内容拥有合法权利，且不侵犯第三方知识产权或隐私。
        </p>

        <h2>五、知识产权</h2>
        <p>
          本网站的设计、标识、界面、算法及除用户内容外的全部内容，其知识产权归我们或相关权利人所有。未经书面许可，不得复制、修改、传播或用于商业目的。
        </p>

        <h2>六、隐私与数据</h2>
        <p>
          我们重视用户隐私。个人信息的收集与使用将遵循我们的隐私政策。我们会采取合理措施保护您的数据安全，但无法保证绝对安全，您理解并接受使用互联网服务的相关风险。
        </p>

        <h2>七、免责声明</h2>
        <p>
          本网站及推荐结果仅供参考，不构成任何专业建议。我们不对因使用或无法使用本网站而产生的直接或间接损失承担责任，法律另有规定的除外。
        </p>

        <h2>八、条款变更</h2>
        <p>
          我们可能适时修订本服务条款，修订后的条款将在本页面公布。若您继续使用本网站，即视为接受修订后的条款。重大变更时，我们可能通过站内通知或注册邮箱提醒您。
        </p>

        <h2>九、适用法律与争议解决</h2>
        <p>
          本条款的订立、履行与解释均适用中华人民共和国法律。因本条款产生的争议，双方应友好协商；协商不成的，可向本网站运营方所在地有管辖权的人民法院提起诉讼。
        </p>

        <h2>十、联系我们</h2>
        <p>
          如对本服务条款有任何疑问，请通过本网站提供的联系方式与我们联系。
        </p>
      </div>

      <router-link to="/" class="terms-back">返回首页</router-link>
    </div>
  </section>
</template>

<script setup>
</script>

<style scoped>
.view-terms {
  padding: 24px 0 48px;
  max-width: 720px;
  margin: 0 auto;
}

.page-title {
  font-family: var(--font-serif);
  font-weight: 600;
  font-size: 24px;
  margin: 0 0 20px;
  color: var(--color-text);
}

.terms-card {
  padding: 28px 24px 24px;
}

.terms-updated {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 0 0 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--color-border);
}

.terms-content {
  font-size: 15px;
  line-height: 1.7;
  color: var(--color-text);
}

.terms-content h2 {
  font-family: var(--font-sans);
  font-weight: 600;
  font-size: 16px;
  margin: 24px 0 10px;
  color: var(--color-text);
}

.terms-content h2:first-child {
  margin-top: 0;
}

.terms-content p {
  margin: 0 0 12px;
}

.terms-content ul {
  margin: 8px 0 12px;
  padding-left: 22px;
}

.terms-content li {
  margin-bottom: 6px;
}

.terms-back {
  display: inline-block;
  margin-top: 28px;
  font-size: 14px;
  color: var(--color-accent);
  text-decoration: none;
}

.terms-back:hover {
  text-decoration: underline;
}
</style>
