import { writable, Writable } from 'svelte/store';

// ginstance is converted from backend config for the GUI
type ginstance = {
  Name: string;
  ImgFolders: string[];
  TimeToQueue: {num: number, unit: string};
  PostInterval: {num: number, unit: string};
  PostDelayAtStartup: string;
  Platforms: {facebook:string, twitter: string}; 
  ItemsInQueue: number;
  NextPostTime: string;
  Status?: string;
}

// export const instances: Writable<instance[]> = writable([])
export const ginstances: Writable<ginstance[]> = writable([])
export const activeInstance: Writable<number> = writable(-1) 
export const unsavedChanges: Writable<boolean> = writable(false)
