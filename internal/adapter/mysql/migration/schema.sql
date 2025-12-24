--
-- Структура таблицы report_calendar
--
CREATE TABLE `report_calendar` (
  `id` varchar(36) NOT NULL,
  `day` int NOT NULL,
  `month` int NOT NULL,
  `year` int NOT NULL,
  `description` varchar(100) DEFAULT NULL,
  `is_paid_vacation` tinyint(1) NOT NULL DEFAULT '1',
  `type_id` varchar(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- --------------------------------------------------------
--
-- Структура таблицы report_standard
--
CREATE TABLE report_standard (
  id varchar(36) NOT NULL,
  month int NOT NULL,
  year int NOT NULL,
  hours int NOT NULL,
  gender_id int NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- --------------------------------------------------------
--
-- Структура таблицы report_type
--
CREATE TABLE report_type (
  id varchar(36) NOT NULL,
  name varchar(50) NOT NULL,
  system_name varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- --------------------------------------------------------
--
-- Структура таблицы report_user
--
CREATE TABLE report_user (
  id varchar(36) NOT NULL,
  user_id varchar(36) NOT NULL,
  day int NOT NULL,
  month int NOT NULL,
  year int NOT NULL,
  hours float NOT NULL DEFAULT 0.0,
  type_id varchar(36) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- --------------------------------------------------------
--
-- Структура таблицы report_vacation
--
CREATE TABLE report_vacation (
  id varchar(36) NOT NULL,
  user_id varchar(36) NOT NULL,
  start_date date NOT NULL,
  end_date date NOT NULL,
  year int NOT NULL,
  description varchar(100),
  status enum('consideration','rejected','approved','') NOT NULL DEFAULT 'consideration',
  create_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- --------------------------------------------------------
--
-- Структура таблицы `report_setting`
--
CREATE TABLE `report_setting` (
  `id` int NOT NULL,
  `vacation_duration` int NOT NULL DEFAULT '30'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
-- --------------------------------------------------------