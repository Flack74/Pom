package web

func getWebUI() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🍅 Pom - Pomodoro Timer</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #0B0F1A 0%, #1a1a2e 100%);
            color: #fff; min-height: 100vh; padding: 20px;
        }
        .container { max-width: 800px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 40px; }
        .header h1 { font-size: 3rem; margin-bottom: 10px; color: #18FFFF; }
        .tabs { display: flex; margin-bottom: 30px; border-radius: 10px; overflow: hidden; }
        .tab { flex: 1; padding: 15px; background: #1a1a2e; border: none; color: #fff; cursor: pointer; transition: all 0.3s; }
        .tab.active { background: #18FFFF; color: #0B0F1A; }
        .tab:hover { background: #FF4081; }
        .content { background: rgba(255,255,255,0.1); border-radius: 15px; padding: 30px; backdrop-filter: blur(10px); }
        .timer-display { text-align: center; margin: 30px 0; }
        .time { font-size: 4rem; font-weight: bold; color: #18FFFF; margin-bottom: 20px; }
        .progress { width: 100%; height: 20px; background: rgba(255,255,255,0.2); border-radius: 10px; overflow: hidden; margin-bottom: 20px; }
        .progress-bar { height: 100%; background: linear-gradient(90deg, #18FFFF, #FF4081); transition: width 0.3s; }
        .controls { display: flex; gap: 15px; justify-content: center; margin: 30px 0; }
        .btn { padding: 12px 24px; border: none; border-radius: 8px; font-size: 1rem; cursor: pointer; transition: all 0.3s; }
        .btn-primary { background: #18FFFF; color: #0B0F1A; }
        .btn-secondary { background: #FF4081; color: #fff; }
        .btn:hover { transform: translateY(-2px); box-shadow: 0 5px 15px rgba(0,0,0,0.3); }
        .settings { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .setting { display: flex; flex-direction: column; }
        .setting label { margin-bottom: 8px; color: #18FFFF; }
        .setting input, .setting select { padding: 10px; border: none; border-radius: 5px; background: rgba(255,255,255,0.1); color: #fff; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 20px; }
        .stat { text-align: center; padding: 20px; background: rgba(255,255,255,0.1); border-radius: 10px; }
        .stat-value { font-size: 2rem; font-weight: bold; color: #18FFFF; }
        .stat-label { color: #ccc; margin-top: 5px; }
        .hidden { display: none; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🍅 Galactic Pomodoro</h1>
            <p>Focus • Break • Achieve</p>
        </div>

        <div class="tabs">
            <button class="tab active" onclick="showTab('timer')">Timer</button>
            <button class="tab" onclick="showTab('dashboard')">Dashboard</button>
        </div>

        <div id="timer-tab" class="content">
            <div class="settings">
                <div class="setting">
                    <label>Profile</label>
                    <select id="profile">
                        <option value="default">Default (25/5)</option>
                        <option value="work">Work (45/10)</option>
                        <option value="study">Study (30/5)</option>
                        <option value="quick">Quick (15/3)</option>
                    </select>
                </div>
                <div class="setting">
                    <label>Work Time (min)</label>
                    <input type="number" id="workTime" value="25" min="1" max="120">
                </div>
                <div class="setting">
                    <label>Break Time (min)</label>
                    <input type="number" id="breakTime" value="5" min="1" max="30">
                </div>
                <div class="setting">
                    <label>Sessions</label>
                    <input type="number" id="sessions" value="4" min="1" max="10">
                </div>
            </div>

            <div class="timer-display">
                <div class="time" id="timeDisplay">25:00</div>
                <div class="progress">
                    <div class="progress-bar" id="progressBar" style="width: 0%"></div>
                </div>
                <p id="sessionInfo">Ready to start • Session 1/4</p>
            </div>

            <div class="controls">
                <button class="btn btn-primary" id="startBtn" onclick="startTimer()">Start Focus</button>
                <button class="btn btn-secondary hidden" id="pauseBtn" onclick="pauseTimer()">Pause</button>
                <button class="btn btn-secondary hidden" id="resumeBtn" onclick="resumeTimer()">Resume</button>
                <button class="btn btn-secondary" id="stopBtn" onclick="stopTimer()">Stop</button>
            </div>
        </div>

        <div id="dashboard-tab" class="content hidden">
            <div class="stats">
                <div class="stat">
                    <div class="stat-value" id="todaySessions">0</div>
                    <div class="stat-label">Today's Sessions</div>
                </div>
                <div class="stat">
                    <div class="stat-value" id="todayMinutes">0</div>
                    <div class="stat-label">Minutes Focused</div>
                </div>
                <div class="stat">
                    <div class="stat-value" id="currentStreak">0</div>
                    <div class="stat-label">Current Streak</div>
                </div>
                <div class="stat">
                    <div class="stat-value" id="totalSessions">0</div>
                    <div class="stat-label">Total Sessions</div>
                </div>
            </div>
        </div>
    </div>

    <script>
        let timer = null;
        let timeLeft = 25 * 60;
        let totalTime = 25 * 60;
        let isRunning = false;
        let currentSession = 1;
        let totalSessions = 4;
        let isBreak = false;

        function showTab(tab) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.content').forEach(c => c.classList.add('hidden'));
            
            event.target.classList.add('active');
            document.getElementById(tab + '-tab').classList.remove('hidden');
            
            if (tab === 'dashboard') loadStats();
        }

        function updateDisplay() {
            const minutes = Math.floor(timeLeft / 60);
            const seconds = timeLeft % 60;
            document.getElementById('timeDisplay').textContent = 
                ` + "`" + `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}` + "`" + `;
            
            const progress = ((totalTime - timeLeft) / totalTime) * 100;
            document.getElementById('progressBar').style.width = progress + '%';
            
            const phase = isBreak ? 'Break' : 'Focus';
            document.getElementById('sessionInfo').textContent = 
                ` + "`" + `${phase} Time • Session ${currentSession}/${totalSessions}` + "`" + `;
        }

        function startTimer() {
            const workTime = parseInt(document.getElementById('workTime').value);
            const breakTime = parseInt(document.getElementById('breakTime').value);
            totalSessions = parseInt(document.getElementById('sessions').value);
            
            timeLeft = isBreak ? breakTime * 60 : workTime * 60;
            totalTime = timeLeft;
            isRunning = true;
            
            document.getElementById('startBtn').classList.add('hidden');
            document.getElementById('pauseBtn').classList.remove('hidden');
            
            timer = setInterval(() => {
                timeLeft--;
                updateDisplay();
                
                if (timeLeft <= 0) {
                    clearInterval(timer);
                    completeSession();
                }
            }, 1000);
        }

        function pauseTimer() {
            clearInterval(timer);
            isRunning = false;
            document.getElementById('pauseBtn').classList.add('hidden');
            document.getElementById('resumeBtn').classList.remove('hidden');
        }

        function resumeTimer() {
            isRunning = true;
            document.getElementById('resumeBtn').classList.add('hidden');
            document.getElementById('pauseBtn').classList.remove('hidden');
            
            timer = setInterval(() => {
                timeLeft--;
                updateDisplay();
                
                if (timeLeft <= 0) {
                    clearInterval(timer);
                    completeSession();
                }
            }, 1000);
        }

        function stopTimer() {
            clearInterval(timer);
            isRunning = false;
            isBreak = false;
            currentSession = 1;
            
            const workTime = parseInt(document.getElementById('workTime').value);
            timeLeft = workTime * 60;
            totalTime = timeLeft;
            
            document.getElementById('startBtn').classList.remove('hidden');
            document.getElementById('pauseBtn').classList.add('hidden');
            document.getElementById('resumeBtn').classList.add('hidden');
            
            updateDisplay();
        }

        function completeSession() {
            if (!isBreak) {
                if (currentSession < totalSessions) {
                    isBreak = true;
                    alert('🎉 Work session complete! Time for a break.');
                    startTimer();
                } else {
                    alert('🏆 All sessions complete! Great work!');
                    stopTimer();
                }
            } else {
                isBreak = false;
                currentSession++;
                alert('☕ Break over! Ready for the next session?');
                document.getElementById('startBtn').classList.remove('hidden');
                document.getElementById('pauseBtn').classList.add('hidden');
                isRunning = false;
            }
        }

        function loadStats() {
            fetch('/api/insights/today')
                .then(r => r.json())
                .then(data => {
                    document.getElementById('todaySessions').textContent = data.sessions || 0;
                    document.getElementById('todayMinutes').textContent = data.minutes || 0;
                })
                .catch(() => {
                    document.getElementById('todaySessions').textContent = '0';
                    document.getElementById('todayMinutes').textContent = '0';
                });
        }

        document.getElementById('profile').addEventListener('change', function() {
            const profiles = {
                default: {work: 25, break: 5, sessions: 4},
                work: {work: 45, break: 10, sessions: 3},
                study: {work: 30, break: 5, sessions: 4},
                quick: {work: 15, break: 3, sessions: 6}
            };
            
            const profile = profiles[this.value];
            document.getElementById('workTime').value = profile.work;
            document.getElementById('breakTime').value = profile.break;
            document.getElementById('sessions').value = profile.sessions;
            
            if (!isRunning) {
                timeLeft = profile.work * 60;
                totalTime = timeLeft;
                totalSessions = profile.sessions;
                updateDisplay();
            }
        });

        updateDisplay();
    </script>
</body>
</html>`
}