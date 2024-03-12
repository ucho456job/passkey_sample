import { useEffect, useState } from 'react';
import { base64UrlStringToArrayBuffer, arrayBufferToBase64Url } from './util';

const LoginForm = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isCredentialRequestPending, setIsCredentialRequestPending] = useState(false);

  const sendAuthenticatorResponseIfWebauthnAvailable = async () => {
    if (isCredentialRequestPending) return;
    setIsCredentialRequestPending(true);

    console.log('sendAuthenticatorResponseIfWebauthnAvailable')
    try {
      if (!(navigator.credentials &&
        navigator.credentials.create &&
        navigator.credentials.get &&
        window.PublicKeyCredential &&
        /* global PublicKeyCredential */
        await PublicKeyCredential.isConditionalMediationAvailable())) {
        return;
      }

      const challengeResponse = await fetch('/api/auth/login-challenge', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
      });
      const challengeData = await challengeResponse.json();
      console.log('challenge response: ', challengeData);

      challengeData.publicKey.challenge = base64UrlStringToArrayBuffer(challengeData.publicKey.challenge);
      console.log('encoded challenge response: ', challengeData);

      const result = await navigator.credentials.get({ publicKey: challengeData.publicKey, mediation: 'conditional' });
      console.log("navigatore.credentials.get result: ", result);

      if (result) {
        const credentials = {
          id: result.id,
          type: result.type,
          rawId: arrayBufferToBase64Url(result.rawId),
          response: {
            clientDataJSON: arrayBufferToBase64Url(result.response.clientDataJSON),
            authenticatorData: arrayBufferToBase64Url(result.response.authenticatorData),
            signature: arrayBufferToBase64Url(result.response.signature),
            userHandle: arrayBufferToBase64Url(result.response.userHandle),
          },
        };
        console.log('create login request body: ', credentials);

        const loginResponse = await fetch('/passkey-session', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(credentials),
        });

        const loginData = await loginResponse.json();
        console.log('create login response: ', loginData);
        alert('Login successful!')
      }
    } catch (e) {
      console.log(e);
      return null;
    } finally {
      setIsCredentialRequestPending(false);
    }
  };

  useEffect(() => {
    console.log('useEffect');
    sendAuthenticatorResponseIfWebauthnAvailable();
  }, []);

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
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
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
