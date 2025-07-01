export interface Profile {
  name: string;
  work_minutes: number;
  break_minutes: number;
  num_sessions: number;
  description: string;
}

export interface TimerSession {
  id: string;
  work_time: number;
  break_time: number;
  sessions: number;
  current_session: number;
  is_running: boolean;
  is_paused: boolean;
  is_break: boolean;
  time_left: number;
  profile: string;
}

export interface Suggestion {
  type: string;
  message: string;
  work_time?: number;
  break_time?: number;
  sessions?: number;
  confidence: number;
}

export interface TodayStats {
  sessions: number;
  minutes: number;
  hours: number;
}

export interface Plugin {
  name: string;
  description: string;
  enabled: boolean;
  triggers: string[];
}