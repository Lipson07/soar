import React, { useState, useEffect, useRef, useCallback } from "react";
import style from "./CallModal.module.scss";
import {
  FiMic,
  FiMicOff,
  FiVideo,
  FiVideoOff,
  FiPhone,
  FiX,
  FiAlertCircle,
  FiUser,
} from "react-icons/fi";

interface CallModalProps {
  isOpen: boolean;
  type: "audio" | "video";
  caller: {
    id: string;
    username: string;
    avatar_url: string | null;
  };
  roomId: string;
  callId: string;
  isIncoming?: boolean;
  ws: WebSocket | null;
  currentUserId: string;
  onAccept?: () => void;
  onReject?: () => void;
  onEnd?: () => void;
  onClose: () => void;
}

function CallModal({
  isOpen,
  type,
  caller,
  roomId,
  callId,
  isIncoming = false,
  ws,
  currentUserId,
  onAccept,
  onReject,
  onEnd,
  onClose,
}: CallModalProps) {
  const [isMuted, setIsMuted] = useState(false);
  const [isVideoOff, setIsVideoOff] = useState(type !== "video");
  const [callDuration, setCallDuration] = useState(0);
  const [callStatus, setCallStatus] = useState<
    "connecting" | "connected" | "ended" | "error"
  >(isIncoming ? "connecting" : "connecting");
  const [error, setError] = useState<string | null>(null);
  const [devicesUnavailable, setDevicesUnavailable] = useState(false);

  const localVideoRef = useRef<HTMLVideoElement>(null);
  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const peerConnectionRef = useRef<RTCPeerConnection | null>(null);
  const localStreamRef = useRef<MediaStream | null>(null);

  const configuration: RTCConfiguration = {
    iceServers: [
      { urls: "stun:stun.l.google.com:19302" },
      { urls: "stun:stun1.l.google.com:19302" },
    ],
  };

  const getAvatarLetter = (name: string) => {
    if (!name) return "?";
    return name.charAt(0).toUpperCase();
  };

  const getAvatarColor = (name: string) => {
    const colors = [
      "#667eea",
      "#764ba2",
      "#f093fb",
      "#f5576c",
      "#4facfe",
      "#00f2fe",
      "#43e97b",
      "#38f9d7",
      "#fa709a",
      "#fee140",
      "#ff6a88",
      "#ff99ac",
    ];
    const index = name
      .split("")
      .reduce((acc, char) => acc + char.charCodeAt(0), 0);
    return colors[index % colors.length];
  };

  useEffect(() => {
    if (isOpen) {
      checkDevices();
    }

    return () => {
      cleanup();
    };
  }, [isOpen]);

  const checkDevices = async () => {
    try {
      const devices = await navigator.mediaDevices.enumerateDevices();
      const hasAudio = devices.some((d) => d.kind === "audioinput");
      const hasVideo = devices.some((d) => d.kind === "videoinput");

      if (!hasAudio) {
        setError("Микрофон не найден");
        setDevicesUnavailable(true);
        return;
      }

      if (type === "video" && !hasVideo) {
        setError("Камера не найдена");
        setDevicesUnavailable(true);
        return;
      }

      if (isIncoming) {
        setCallStatus("connecting");
      } else {
        startCall();
      }
    } catch (err) {
      console.error("Failed to check devices:", err);
    }
  };

  useEffect(() => {
    if (callStatus === "connected") {
      const timer = setInterval(() => {
        setCallDuration((prev) => prev + 1);
      }, 1000);
      return () => clearInterval(timer);
    }
  }, [callStatus]);

  const cleanup = () => {
    if (peerConnectionRef.current) {
      peerConnectionRef.current.close();
      peerConnectionRef.current = null;
    }
    if (localStreamRef.current) {
      localStreamRef.current.getTracks().forEach((track) => track.stop());
      localStreamRef.current = null;
    }
  };

  const getMediaStream = async (): Promise<MediaStream | null> => {
    const constraints: MediaStreamConstraints = {
      audio: true,
      video: type === "video" ? { width: 1280, height: 720 } : false,
    };

    try {
      return await navigator.mediaDevices.getUserMedia(constraints);
    } catch (err: any) {
      console.error("Failed to get media stream:", err);

      if (err.name === "NotFoundError" || err.name === "DevicesNotFoundError") {
        setError(
          type === "video"
            ? "Камера или микрофон не найдены"
            : "Микрофон не найден",
        );
      } else if (
        err.name === "NotAllowedError" ||
        err.name === "PermissionDeniedError"
      ) {
        setError("Нет доступа к камере или микрофону");
      } else if (err.name === "NotReadableError") {
        setError("Устройство занято другим приложением");
      } else {
        setError(`Ошибка доступа к устройствам: ${err.message}`);
      }

      setDevicesUnavailable(true);
      setCallStatus("error");
      return null;
    }
  };

  const startCall = async () => {
    const stream = await getMediaStream();
    if (!stream) return;

    localStreamRef.current = stream;
    if (localVideoRef.current) {
      localVideoRef.current.srcObject = stream;
    }

    try {
      await createPeerConnection(stream);

      const offer = await peerConnectionRef.current!.createOffer();
      await peerConnectionRef.current!.setLocalDescription(offer);

      sendSignal({
        type: "offer",
        room_id: roomId,
        call_id: callId,
        sdp: offer.sdp,
      });

      setCallStatus("connected");
    } catch (error) {
      console.error("Failed to start call:", error);
      setError("Не удалось установить соединение");
      setCallStatus("error");
    }
  };

  const acceptCall = async () => {
    const stream = await getMediaStream();
    if (!stream) {
      rejectCall();
      return;
    }

    localStreamRef.current = stream;
    if (localVideoRef.current) {
      localVideoRef.current.srcObject = stream;
    }

    try {
      await createPeerConnection(stream);

      sendSignal({
        type: "call-accept",
        room_id: roomId,
        call_id: callId,
      });

      setCallStatus("connected");
      onAccept?.();
    } catch (error) {
      console.error("Failed to accept call:", error);
      rejectCall();
    }
  };

  const createPeerConnection = async (stream: MediaStream) => {
    const pc = new RTCPeerConnection(configuration);
    peerConnectionRef.current = pc;

    stream.getTracks().forEach((track) => {
      pc.addTrack(track, stream);
    });

    pc.ontrack = (event) => {
      if (remoteVideoRef.current) {
        remoteVideoRef.current.srcObject = event.streams[0];
      }
    };

    pc.onicecandidate = (event) => {
      if (event.candidate) {
        sendSignal({
          type: "ice-candidate",
          room_id: roomId,
          call_id: callId,
          candidate: event.candidate,
        });
      }
    };

    pc.onconnectionstatechange = () => {
      if (
        pc.connectionState === "disconnected" ||
        pc.connectionState === "failed"
      ) {
        endCall();
      }
    };
  };

  const handleSignal = useCallback(
    async (signal: any) => {
      const pc = peerConnectionRef.current;
      if (!pc) return;

      try {
        if (signal.type === "offer") {
          await pc.setRemoteDescription({ type: "offer", sdp: signal.sdp });
          const answer = await pc.createAnswer();
          await pc.setLocalDescription(answer);

          sendSignal({
            type: "answer",
            room_id: roomId,
            call_id: callId,
            sdp: answer.sdp,
          });
        } else if (signal.type === "answer") {
          await pc.setRemoteDescription({ type: "answer", sdp: signal.sdp });
        } else if (signal.type === "ice-candidate") {
          await pc.addIceCandidate(new RTCIceCandidate(signal.candidate));
        }
      } catch (error) {
        console.error("Error handling signal:", error);
      }
    },
    [roomId, callId],
  );

  useEffect(() => {
    if (ws) {
      const originalOnMessage = ws.onmessage;
      ws.onmessage = (event) => {
        const signal = JSON.parse(event.data);
        if (signal.room_id === roomId) {
          if (signal.type === "call-accept") {
            setCallStatus("connected");
          } else if (signal.type === "call-reject") {
            onClose();
          } else if (signal.type === "call-end") {
            onClose();
          } else {
            handleSignal(signal);
          }
        }
        if (originalOnMessage) {
          originalOnMessage.call(ws, event);
        }
      };
    }
  }, [ws, roomId, handleSignal, onClose]);

  const sendSignal = (signal: any) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(signal));
    }
  };

  const toggleMute = () => {
    if (localStreamRef.current) {
      localStreamRef.current.getAudioTracks().forEach((track) => {
        track.enabled = !track.enabled;
      });
      setIsMuted(!isMuted);
    }
  };

  const toggleVideo = () => {
    if (localStreamRef.current) {
      localStreamRef.current.getVideoTracks().forEach((track) => {
        track.enabled = !track.enabled;
      });
      setIsVideoOff(!isVideoOff);
    }
  };

  const endCall = () => {
    sendSignal({
      type: "call-end",
      room_id: roomId,
      call_id: callId,
    });
    cleanup();
    onEnd?.();
    onClose();
  };

  const rejectCall = () => {
    sendSignal({
      type: "call-reject",
      room_id: roomId,
      call_id: callId,
    });
    onReject?.();
    onClose();
  };

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  if (!isOpen) return null;

  return (
    <div className={style.overlay}>
      <div className={style.callModal}>
        {/* Фон с размытием */}
        <div className={style.background}>
          <div className={style.bgGradient} />
          <div className={style.bgPattern} />
        </div>

        {/* Контент */}
        <div className={style.content}>
          {/* Аватар и информация */}
          <div className={style.callerSection}>
            <div className={style.avatarWrapper}>
              {caller.avatar_url ? (
                <img
                  src={caller.avatar_url}
                  alt={caller.username}
                  className={style.avatar}
                />
              ) : (
                <div
                  className={style.avatarPlaceholder}
                  style={{ background: getAvatarColor(caller.username) }}
                >
                  <span>{getAvatarLetter(caller.username)}</span>
                </div>
              )}
              {callStatus === "connected" && type === "audio" && (
                <div className={style.avatarPulse} />
              )}
            </div>

            <h2 className={style.callerName}>{caller.username}</h2>

            <div className={style.callStatus}>
              {callStatus === "connecting" && (
                <>
                  <div className={style.spinnerSmall} />
                  <span>
                    {isIncoming ? "Входящий звонок..." : "Соединение..."}
                  </span>
                </>
              )}
              {callStatus === "connected" && (
                <span>{formatDuration(callDuration)}</span>
              )}
              {callStatus === "error" && (
                <span className={style.errorText}>{error}</span>
              )}
            </div>
          </div>

          {/* Видео */}
          {type === "video" && callStatus !== "error" && (
            <>
              <div className={style.remoteVideo}>
                <video ref={remoteVideoRef} autoPlay playsInline />
              </div>
              <div className={style.localVideo}>
                <video ref={localVideoRef} autoPlay playsInline muted />
              </div>
            </>
          )}

          {/* Ошибка */}
          {callStatus === "error" && (
            <div className={style.errorBlock}>
              <FiAlertCircle size={48} />
              <p>{error || "Ошибка соединения"}</p>
            </div>
          )}

          {/* Кнопки управления */}
          <div className={style.controls}>
            {callStatus === "connecting" &&
            isIncoming &&
            !devicesUnavailable ? (
              <>
                <button className={style.acceptBtn} onClick={acceptCall}>
                  <FiPhone size={28} />
                </button>
                <button className={style.rejectBtn} onClick={rejectCall}>
                  <FiX size={28} />
                </button>
              </>
            ) : callStatus === "error" || devicesUnavailable ? (
              <button className={style.closeBtn} onClick={onClose}>
                <FiX size={24} />
                <span>Закрыть</span>
              </button>
            ) : (
              <>
                <button
                  className={`${style.controlBtn} ${isMuted ? style.active : ""}`}
                  onClick={toggleMute}
                >
                  {isMuted ? <FiMicOff size={22} /> : <FiMic size={22} />}
                </button>
                {type === "video" && (
                  <button
                    className={`${style.controlBtn} ${isVideoOff ? style.active : ""}`}
                    onClick={toggleVideo}
                  >
                    {isVideoOff ? (
                      <FiVideoOff size={22} />
                    ) : (
                      <FiVideo size={22} />
                    )}
                  </button>
                )}
                <button className={style.endBtn} onClick={endCall}>
                  <FiPhone size={28} />
                </button>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default CallModal;
