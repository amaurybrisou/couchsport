<script>
  import { AUTH_REQUEST } from 'actions/auth'
  import AuthForm from 'components/auth/AuthForm'

  export default {
    name: 'Login',
    functional: true,
    props: { welcome: { type: String, default: '' } },
    render(createElement, context) {
      context.props = {
        submit: function (user) {
          context.parent.$store
            .dispatch(AUTH_REQUEST, user)
            .then(() => {
              context.parent.$router.push({ name: 'profile' })
            })
            .catch((data) => {
              console.log(data)
              this.errors = [data]
            })
        },
        welcome: context.props.welcome,
        title: context.parent.$t('login'),
        buttonMessage: context.parent.$t('login')
      }

      return createElement(AuthForm, context)
    }
  }
</script>
