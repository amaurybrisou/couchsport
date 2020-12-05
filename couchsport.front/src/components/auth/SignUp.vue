<script>
  import { repositoryFactory } from 'repos/repositoryFactory'

  import AuthForm from 'components/auth/AuthForm'
  const userRepository = repositoryFactory.get('user')

  export default {
    name: 'SignUp',
    functional: true,
    render(createElement, context) {
      context.props = {
        submit: function (user) {
          userRepository
            .create(user)
            .then(() => {
              context.parent.$router.push({
                name: 'login',
                params: {
                  welcome: context.parent.$t('message.signup_success_welcome')
                }
              })
            })
            .catch(({ response: { data } }) => {
              this.errors = [data]
            })
        },
        title: context.parent.$t('signup'),
        buttonMessage: context.parent.$t('signup')
      }
      return createElement(AuthForm, context)
    }
  }
</script>
