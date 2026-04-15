<script lang="ts">
  import { onMount } from "svelte";
  import logoUrl from "./assets/logo.png";
  import { GetConfig, SelectGtaFolder, SelectReduxFolder, SetActiveRedux, ToggleAutoInject, RemoveReduxFolder, SetRunOnStartup, SetStartWindowBox } from "../wailsjs/go/main/App";
  import { EventsOn } from "../wailsjs/runtime/runtime";
  import type { main } from "../wailsjs/go/models";

  let config: main.Config = {
    gta_path: "",
    redux_folders: [],
    active_redux: "",
    auto_inject: false,
    run_on_startup: false,
    start_window_box: true
  };

  let statusMessage = "Ожидание настройки...";
  let currentTab = "dashboard"; // "dashboard" | "settings"

  async function loadConfig() {
    try {
      config = await GetConfig();
    } catch (err) {
      console.error(err);
    }
  }

  onMount(() => {
    loadConfig();
    EventsOn("status_update", (msg: string) => {
      statusMessage = msg;
    });
  });

  async function chooseGtaPath() { await SelectGtaFolder(); await loadConfig(); }
  async function addRedux() { await SelectReduxFolder(); await loadConfig(); }
  async function removeRedux(path: string) { await RemoveReduxFolder(path); await loadConfig(); }
  async function selectRedux(path: string) { await SetActiveRedux(path); await loadConfig(); }
  async function toggleAutoInject() { await ToggleAutoInject(!config.auto_inject); await loadConfig(); }
  async function toggleStartup() { await SetRunOnStartup(!config.run_on_startup); await loadConfig(); }
  async function toggleShowBox() { await SetStartWindowBox(!config.start_window_box); await loadConfig(); }
</script>

<main class="app-container">
  <!-- Sidebar -->
  <aside class="sidebar">
    <div class="logo">
      <img src={logoUrl} alt="FasterRedux Logo" class="brand-icon" />
      <h1>Faster<span>Redux</span></h1>
    </div>

    <nav class="nav-menu">
      <button class:active={currentTab === 'dashboard'} on:click={() => currentTab = 'dashboard'}>Мои Редуксы</button>
      <button class:active={currentTab === 'settings'} on:click={() => currentTab = 'settings'}>Настройки Приложения</button>
    </nav>
  </aside>

  <!-- Main Content -->
  <section class="content">
    <header class="header">
      <div class="status-panel" class:active={config.auto_inject}>
        <div class="status-info">
          <span>Статус защиты</span>
          <h2>{statusMessage}</h2>
        </div>
        <button class="toggle-btn" class:on={config.auto_inject} on:click={toggleAutoInject}>
          {config.auto_inject ? "Остановить автообход" : "Запустить автообход"}
        </button>
      </div>
    </header>

    <div class="scroll-container">
      {#if currentTab === 'dashboard'}
        <div class="list-header">
          <h2>Библиотека Редуксов</h2>
          <button class="add-btn" on:click={addRedux}>+ Добавить Редукс</button>
        </div>

        {#if config.redux_folders && config.redux_folders.length > 0}
          <div class="grid">
            {#each config.redux_folders as redux}
              <div class="redux-card" class:selected={config.active_redux === redux}>
                <div class="card-info">
                  <h3>{redux.split('\\').pop() || redux.split('/').pop()}</h3>
                  <p title={redux}>{redux}</p>
                </div>
                <div class="card-actions">
                  {#if config.active_redux !== redux}
                    <button class="apply-btn" on:click={() => selectRedux(redux)}>Активировать</button>
                  {:else}
                    <span class="active-badge">Активен</span>
                  {/if}
                  <button class="del-btn" on:click={() => removeRedux(redux)}>Удалить</button>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <div class="empty-state">
            <h3>Здесь пока пусто</h3>
            <p>Добавьте хотя бы один редукс, чтобы начать.</p>
          </div>
        {/if}
      {:else}
        <!-- Settings Tab -->
        <div class="settings-container">
          <h2>Настройки FasterRedux</h2>

          <div class="setting-group">
            <div class="setting-item">
              <div class="setting-text">
                <label>Путь к игре (GTA 5)</label>
                <span>Укажите папку, где лежит GTA5.exe</span>
              </div>
              <div class="path-box">
                <input type="text" readonly value={config.gta_path || 'Не выбран...'} />
                <button on:click={chooseGtaPath}>Обзор</button>
              </div>
            </div>

            <div class="setting-item">
              <div class="setting-text">
                <label>Запускать вместе с Windows</label>
                <span>Программа будет автоматически стартовать при включении ПК</span>
              </div>
              <label class="switch">
                <input type="checkbox" checked={config.run_on_startup} on:change={toggleStartup}>
                <span class="slider"></span>
              </label>
            </div>

            <div class="setting-item">
              <div class="setting-text">
                <label>Показывать окно при запуске</label>
                <span>Если выключить, программа будет стартовать в фоне (без интерфейса)</span>
              </div>
              <label class="switch">
                <input type="checkbox" checked={config.start_window_box} on:change={toggleShowBox}>
                <span class="slider"></span>
              </label>
            </div>
            
            <div class="tip-box">
              <strong>Подсказка:</strong> При включении свойства "Запустить вместе с Windows" Автообход также будет автоматически включен.<br>
              При закрытии окна (крестик) программа свернется в <b>системный трей</b> (возле часов) и останется работать в фоновом режиме для обхода блокировки.
            </div>
          </div>
        </div>
      {/if}
    </div>
  </section>
</main>

<style>
  :global(body) {
    margin: 0; padding: 0;
    font-family: 'Segoe UI', system-ui, sans-serif;
    color: #e2e8f0; background: #0f111a;
    user-select: none;
  }
  .app-container { display: flex; height: 100vh; overflow: hidden; }

  /* Sidebar */
  .sidebar { width: 260px; background: #151722; padding: 25px 0; border-right: 1px solid #1f2233; display: flex; flex-direction: column; }
  .logo { padding: 0 25px 30px; display: flex; align-items: center; gap: 12px; }
  .brand-icon { width: 45px; height: 45px; filter: drop-shadow(0 4px 8px rgba(106, 53, 255, 0.4)); }
  .logo h1 { margin: 0; font-size: 24px; color: #fff; font-weight: 800; letter-spacing: 1px; }
  .logo span { color: #6A35FF; }

  .nav-menu { display: flex; flex-direction: column; gap: 5px; padding: 0 15px; }
  .nav-menu button {
    background: transparent; color: #94a3b8; border: none; text-align: left; padding: 12px 20px; font-size: 15px; border-radius: 8px; cursor: pointer; transition: 0.2s; font-weight: 500;
  }
  .nav-menu button:hover { background: #1e2235; color: #fff; }
  .nav-menu button.active { background: #6A35FF; color: #fff; box-shadow: 0 4px 12px rgba(106, 53, 255, 0.3); }

  /* Main Content */
  .content { flex: 1; display: flex; flex-direction: column; }
  .header { padding: 20px 30px; background: #11131c; border-bottom: 1px solid #1f2233; }

  .status-panel { display: flex; justify-content: space-between; align-items: center; background: #1a1d2d; padding: 15px 25px; border-radius: 12px; border-left: 5px solid #475569; transition: 0.3s; box-shadow: 0 4px 6px rgba(0,0,0,0.2); }  
  .status-panel.active { border-left-color: #6A35FF; background: rgba(106, 53, 255, 0.1); }     
  .status-info span { font-size: 12px; text-transform: uppercase; color: #94a3b8; font-weight: bold; letter-spacing: 1px; }
  .status-info h2 { margin: 5px 0 0; font-size: 18px; color: #fff; }

  .toggle-btn { padding: 10px 24px; border-radius: 8px; border: none; font-weight: bold; background: #334155; color: #fff; cursor: pointer; transition: 0.2s; font-size: 14px; }
  .toggle-btn.on { background: #ef4444; box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3); }
  .toggle-btn:not(.on):hover { background: #6A35FF; box-shadow: 0 4px 12px rgba(106, 53, 255, 0.3); }

  .scroll-container { padding: 30px; flex: 1; overflow-y: auto; overflow-x: hidden; }
  .list-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 25px; }
  .list-header h2 { margin: 0; font-size: 22px; font-weight: 600; color: #fff; }
  .add-btn { background: #6A35FF; color: #fff; border: none; padding: 10px 20px; border-radius: 8px; cursor: pointer; font-weight: bold; transition: 0.2s; box-shadow: 0 4px 10px rgba(106, 53, 255, 0.3); }    
  .add-btn:hover { background: #5824e6; box-shadow: 0 6px 15px rgba(106, 53, 255, 0.4); transform: translateY(-1px); }

  .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 20px; }
  .redux-card { background: #1a1d2d; border: 1px solid #2b3046; border-radius: 12px; padding: 20px; display: flex; flex-direction: column; justify-content: space-between; gap: 20px; transition: 0.2s; }
  .redux-card:hover { border-color: #4a5578; transform: translateY(-2px); }
  .redux-card.selected { border-color: #6A35FF; background: rgba(106, 53, 255, 0.05); box-shadow: 0 0 20px rgba(106, 53, 255, 0.15); }
  .card-info h3 { margin: 0 0 8px; font-size: 18px; color: #f8fafc; }
  .card-info p { margin: 0; font-size: 13px; color: #94a3b8; overflow: hidden; text-overflow: ellipsis; }

  .card-actions { display: flex; justify-content: flex-end; align-items: center; gap: 10px; }
  .active-badge { font-size: 12px; color: #a78bfa; font-weight: bold; text-transform: uppercase; margin-right: auto; letter-spacing: 0.5px; }
  .apply-btn { background: #6A35FF; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; font-size: 13px; font-weight: 600; transition: 0.2s; }    
  .apply-btn:hover { background: #5824e6; }
  .del-btn { background: transparent; color: #ef4444; border: 1px solid #ef4444; padding: 7px 15px; border-radius: 6px; cursor: pointer; font-size: 13px; font-weight: 600; transition: 0.2s; }
  .del-btn:hover { background: #ef4444; color: #fff; }

  .empty-state { text-align: center; padding: 60px 0; background: #1a1d2d; border-radius: 12px; border: 2px dashed #4a5578; }
  .empty-state h3 { margin: 0 0 10px; font-size: 20px; color: #e2e8f0; }        
  .empty-state p { margin: 0; color: #94a3b8; }

  /* Settings UI */
  .settings-container h2 { margin: 0 0 25px; font-size: 22px; color: #fff; }
  .setting-group { display: flex; flex-direction: column; gap: 25px; background: #1a1d2d; padding: 30px; border-radius: 12px; border: 1px solid #2b3046; }      
  .setting-item { display: flex; justify-content: space-between; align-items: center; }
  .setting-text label { display: block; font-size: 16px; font-weight: 600; color: #f1f5f9; margin-bottom: 4px; }
  .setting-text span { font-size: 13px; color: #94a3b8; }

  .path-box { display: flex; gap: 10px; width: 350px; }
  .path-box input { flex: 1; background: #0f111a; border: 1px solid #2b3046; color: #f1f5f9; padding: 10px 15px; border-radius: 8px; font-size: 13px; outline: none; transition: 0.2s; }
  .path-box input:focus { border-color: #6A35FF; box-shadow: 0 0 0 2px rgba(106, 53, 255, 0.2); }
  .path-box button { background: #6A35FF; color: white; border: none; border-radius: 8px; padding: 0 15px; cursor: pointer; font-weight: 600; transition: 0.2s; }
  .path-box button:hover { background: #5824e6; }

  /* iOS Switch */
  .switch { position: relative; display: inline-block; width: 44px; height: 24px; }
  .switch input { opacity: 0; width: 0; height: 0; }
  .slider { position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0; background-color: #334155; transition: .4s; border-radius: 24px; }       
  .slider:before { position: absolute; content: ""; height: 18px; width: 18px; left: 3px; bottom: 3px; background-color: white; transition: .4s; border-radius: 50%; box-shadow: 0 2px 4px rgba(0,0,0,0.2); }
  input:checked + .slider { background-color: #6A35FF; box-shadow: 0 0 10px rgba(106, 53, 255, 0.4); }
  input:checked + .slider:before { transform: translateX(20px); }

  .tip-box { margin-top: 15px; padding: 15px; background: rgba(106, 53, 255, 0.1); border-left: 4px solid #6A35FF; border-radius: 8px; font-size: 14px; color: #d8b4fe; line-height: 1.5; }

</style>
