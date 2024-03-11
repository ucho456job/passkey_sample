import { useState } from 'react';

const LoginForm = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleEmailFocus = async () => {
    try {
      // if (!(navigator.credentials &&
      //   navigator.credentials.create &&
      //   navigator.credentials.get &&
      //   window.PublicKeyCredential &&
      //   await PublicKeyCredential.isConditionalMediationAvailable())) {
      //   return;
      // }
  
      // const optionsJSON = await backend.fetchWebauthnAssertionOptions();
      // if (optionsJSON != null) {
      //   options = webauthn.parseRequestOptionsFromJSON(optionsJSON);
      // } else {
      //   return null;
      // }
  
      // options['mediation'] = 'conditional';
  
      // const response = await navigator.credentials.get(options);
      // return await backend.postWebauthnAssertion(response.toJSON()); // send the authenticator response to the backend
    } catch (e) {
      console.log(e);
      return null;
    }
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    console.log('Login submitted for email:', email);
  };

  return (
    <div>
      <h1>Login Form</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="email">Email:</label>
          <input
            id="email"
            type="text"
            name="username"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            onFocus={handleEmailFocus}
            autocomplete="username webauthn"
            required
          />
        </div>
        <div>
          <label htmlFor="password">Password:</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <div>
          <button type="submit">Login</button>
        </div>
      </form>
    </div>
  );
};

export default LoginForm;
