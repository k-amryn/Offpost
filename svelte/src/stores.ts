import { writable, Writable } from 'svelte/store';

// ginstance is converted from backend config for the GUI
type ginstance = {
  Name: string;
  ImgFolders: string[];
  QueueDelay: {num: number, unit: string};
  PostDelay: {num: number, unit: string};
  StartupPostDelay: string;
  Platforms: {}; 
  Caption: string;
  ItemsInQueue: number;
  NextPostTime: string;
  Status?: string;
  Image: string;
}

// export const instances: Writable<instance[]> = writable([])
export const ginstances: Writable<ginstance[]> = writable([])
export const ginstancesOld: Writable<ginstance[]> = writable([])

export const activeInstance: Writable<number> = writable(-1) 
export const unsavedChanges: Writable<boolean> = writable(false)
