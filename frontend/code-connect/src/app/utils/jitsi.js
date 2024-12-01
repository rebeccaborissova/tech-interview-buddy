export const generateJitsiRoom = () => {
  const roomName = `room-${Math.random().toString(36).substring(2, 15)}`;
  const jitsiDomain = "meet.jit.si";
  return `https://${jitsiDomain}/${roomName}`;
};