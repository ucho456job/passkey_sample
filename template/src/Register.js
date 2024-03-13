import { base64UrlStringToArrayBuffer, arrayBufferToBase64Url } from './util';

export default function Register() {
  const handleGeneratePasskey = async () => {
    console.log('handleGeneratePasskey')

    // Check if the browser supports WebAuthn
    if (!(navigator.credentials &&
      navigator.credentials.create &&
      navigator.credentials.get &&
      window.PublicKeyCredential &&
      /* global PublicKeyCredential */
      await PublicKeyCredential.isConditionalMediationAvailable())) {
      return;
    }

    try {
      // Receive a challenge from the server
      const challengeResponse = await fetch('/register_challenge', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
      });
      const challengeData = await challengeResponse.json();
      challengeData.publicKey.challenge = base64UrlStringToArrayBuffer(challengeData.publicKey.challenge);
      challengeData.publicKey.user.id = base64UrlStringToArrayBuffer(challengeData.publicKey.user.id);
      if (challengeData.publicKey.excludeCredentials) {
        challengeData.publicKey.excludeCredentials = challengeData.publicKey.excludeCredentials.map(cred => {
          return {
            ...cred,
            id: base64UrlStringToArrayBuffer(cred.id),
          };
        });
      }
      console.log('encoded challenge response: ', challengeData);

      // Get the authenticator response
      const result = await navigator.credentials.create(challengeData);
      console.log("navigatore.credentials.create result: ", result);

      if (result) {
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
  
        // Send the authenticator response to the server
        const passkeyResponse = await fetch('/passkeys', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(credentials),
        });
  
        const passkeyData = await passkeyResponse.json();
        console.log('create passkey response: ', passkeyData);
        alert('Passkey created successfully!')
      }
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
