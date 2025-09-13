import { setupWorker } from "msw/browser";
import {
  getLoginUserMockHandler,
  getLogoutUserMockHandler,
  getRegisterUserMockHandler,
  getGetProfileMockHandler,
  getUpdateProfileMockHandler,
  getDeleteUserMockHandler,
  getAdminListUsersMockHandler,
  getAdminGetUserMockHandler,
  getAdminUpdateUserMockHandler,
} from "../services/api/user/user.msw";
import { getHealthCheckMockHandler } from "../services/api/health/health.msw";
import {
  getGetGamesMockHandler,
  getCreateGameMockHandler,
  getAdminListGamesMockHandler,
  getUpdateGameMockHandler,
  getDeleteGameMockHandler,
} from "../services/api/games/games.msw";
import {
  getGetPicksMockHandler,
  getSubmitPicksMockHandler,
  getAdminSubmitPicksMockHandler,
  getAdminListPicksMockHandler,
  getAdminGetPicksByWeekMockHandler,
  getAdminGetPicksByUserMockHandler,
  getAdminDeletePickMockHandler,
} from "../services/api/picks/picks.msw";
import {
  getGetWeeklyResultsMockHandler,
  getGetSeasonResultsMockHandler,
  getSubmitResultMockHandler,
} from "../services/api/results/results.msw";
import {
  getGetSurvivorPicksMockHandler,
  getSubmitSurvivorPickMockHandler,
} from "../services/api/survivor/survivor.msw";

export const worker = setupWorker(
  getLoginUserMockHandler(),
  getLogoutUserMockHandler(),
  getRegisterUserMockHandler(),
  getHealthCheckMockHandler(),
  getGetProfileMockHandler(),
  getUpdateProfileMockHandler(),
  getDeleteUserMockHandler(),
  getAdminListUsersMockHandler(),
  getAdminGetUserMockHandler(),
  getAdminUpdateUserMockHandler(),
  getGetGamesMockHandler(),
  getCreateGameMockHandler(),
  getAdminListGamesMockHandler(),
  getUpdateGameMockHandler(),
  getDeleteGameMockHandler(),
  getGetPicksMockHandler(),
  getSubmitPicksMockHandler(),
  getAdminSubmitPicksMockHandler(),
  getAdminListPicksMockHandler(),
  getAdminGetPicksByWeekMockHandler(),
  getAdminGetPicksByUserMockHandler(),
  getAdminDeletePickMockHandler(),
  getGetWeeklyResultsMockHandler(),
  getGetSeasonResultsMockHandler(),
  getSubmitResultMockHandler(),
  getGetSurvivorPicksMockHandler(),
  getSubmitSurvivorPickMockHandler(),
);
