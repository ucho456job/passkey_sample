import { base64UrlStringToArrayBuffer, arrayBufferToBase64Url } from './util';

export default function Register() {
  const handleGeneratePasskey = async () => {
    console.log('handleGeneratePasskey')
    try {
      const challengeResponse = await fetch('/api/auth/register-challenge', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
      });
      const challengeData = await challengeResponse.json();
      console.log('challenge response: ', challengeData);

      challengeData.publicKey.challenge = base64UrlStringToArrayBuffer(challengeData.publicKey.challenge);
      challengeData.publicKey.user.id = base64UrlStringToArrayBuffer(challengeData.publicKey.user.id);

      console.log('encoded challenge response: ', challengeData);

      const result = await navigator.credentials.create(challengeData);
      console.log("navigatore.credentials.create result: ", result);

      const credentials = {
        id: result.id,
        type: result.type,
        rawId: arrayBufferToBase64Url(result.rawId),
        response: {
          clientDataJSON: arrayBufferToBase64Url(result.response.clientDataJSON),
          attestationObject: arrayBufferToBase64Url(result.response.attestationObject),
        },
      };
      console.log('create passkey request body: ', credentials);

      const passkeyResponse = await fetch('/api/auth/passkey', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });

      const passkeyData = await passkeyResponse.json();
      console.log('create passkey response: ', passkeyData);
      alert('Passkey created successfully!')
    } catch (err) {
      console.error('Error:', err);
      alert('An error occurred.');
    }
  };

  return (
    <div>
      <h1>Register</h1>
      <button onClick={handleGeneratePasskey}>Generate Passkey</button>
    </div>
  );
}
