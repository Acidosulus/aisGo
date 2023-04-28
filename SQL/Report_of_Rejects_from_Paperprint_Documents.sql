SELECT fio1 as fio,
	sum(1) as total,
	sum(case when len(lk)=0 or lk is null then 0 else 1 end) as lk,
	sum(case when edo is null then 0 else 1 end) as edo,
	sum(case when reject is null then 0 else 1 end) as reject
from(select agreement.[Номер] as nc,
	staff1.ФИО as fio1,
	staff2.ФИО as fio2,
	staff3.ФИО as fio3,
	staff4.ФИО as fio4,
	reject.[Название] as reject,
	edo.[Название] as edo,
	lk = case when len(lk.nc)=10 then 'ЛК' else '' end
from stack.[Договор] agreement
left join (select 	stack.[Договор].Номер,
					stack.[Виды параметров].Название
		from stack.[Договор], stack.[Свойства]
		left join stack.[Виды параметров] on stack.[Виды параметров].ROW_ID  = stack.[Свойства].[Виды-Параметры] 
		where 	('2023-04-30' between stack.[Свойства].ДатНач and stack.[Свойства].ДатКнц) and
				stack.[Свойства].[Параметры-Договор] = stack.[Договор].ROW_ID and  
				stack.[Виды параметров].Название = 'ОТКАЗ_БДОК') as reject on reject.[Номер] = agreement.Номер
left join (select 	stack.[Договор].Номер,
					stack.[Виды параметров].Название
		from stack.[Договор], stack.[Свойства]
		left join stack.[Виды параметров] on stack.[Виды параметров].ROW_ID  = stack.[Свойства].[Виды-Параметры] 
		where 	('2023-04-30' between stack.[Свойства].ДатНач and stack.[Свойства].ДатКнц) and
				stack.[Свойства].[Параметры-Договор] = stack.[Договор].ROW_ID and  
				stack.[Виды параметров].Название = 'ЭДО') as edo on edo.[Номер] = agreement.Номер
left join (			select distinct left(agr.[Номер],10) as nc
							from stack.[Пароли привязка] pp, stack.[Договор] agr
							where 		agr.row_id = pp.[Привязка-Договор]
										and pp.[Состояние] > 0 ) as lk on lk.nc = agreement.Номер 
left join stack.[Сотрудники] as staff1 on staff1.ROW_ID = agreement.Сотрудник1
left join stack.[Сотрудники] as staff2 on staff2.ROW_ID = agreement.Сотрудник2
left join stack.[Сотрудники] as staff3 on staff3.ROW_ID = agreement.Сотрудник3
left join stack.[Сотрудники] as staff4 on staff4.ROW_ID = agreement.Сотрудник4
WHERE 	('2023-04-01' between agreement.[Начало договора] and agreement.[Окончание]) and 
		('2023-04-30' between agreement.[Начало договора] and agreement.[Окончание]) and 
		LEN( agreement.[Номер]) = 10) as orul
where fio1 is not null
group by fio1
union all 
select '' as fio1, null as total, null as lk, null as edo, null as reject
union all
SELECT fio3 as fio,
	sum(1) as total,
	sum(case when len(lk)=0 or lk is null then 0 else 1 end) as lk,
	sum(case when edo is null then 0 else 1 end) as edo,
	sum(case when reject is null then 0 else 1 end) as reject
from(select agreement.[Номер] as nc,
	staff1.ФИО as fio1,
	staff2.ФИО as fio2,
	staff3.ФИО as fio3,
	staff4.ФИО as fio4,
	reject.[Название] as reject,
	edo.[Название] as edo,
	lk = case when len(lk.nc)=10 then 'ЛК' else '' end
from stack.[Договор] agreement
left join (select 	stack.[Договор].Номер,
					stack.[Виды параметров].Название
		from stack.[Договор], stack.[Свойства]
		left join stack.[Виды параметров] on stack.[Виды параметров].ROW_ID  = stack.[Свойства].[Виды-Параметры] 
		where 	('2023-04-30' between stack.[Свойства].ДатНач and stack.[Свойства].ДатКнц) and
				stack.[Свойства].[Параметры-Договор] = stack.[Договор].ROW_ID and  
				stack.[Виды параметров].Название = 'ОТКАЗ_БДОК') as reject on reject.[Номер] = agreement.Номер
left join (select 	stack.[Договор].Номер,
					stack.[Виды параметров].Название
		from stack.[Договор], stack.[Свойства]
		left join stack.[Виды параметров] on stack.[Виды параметров].ROW_ID  = stack.[Свойства].[Виды-Параметры] 
		where 	('2023-04-30' between stack.[Свойства].ДатНач and stack.[Свойства].ДатКнц) and
				stack.[Свойства].[Параметры-Договор] = stack.[Договор].ROW_ID and  
				stack.[Виды параметров].Название = 'ЭДО') as edo on edo.[Номер] = agreement.Номер
left join (			select distinct left(agr.[Номер],10) as nc
							from stack.[Пароли привязка] pp, stack.[Договор] agr
							where 		agr.row_id = pp.[Привязка-Договор]
										and pp.[Состояние] > 0 ) as lk on lk.nc = agreement.Номер 
left join stack.[Сотрудники] as staff1 on staff1.ROW_ID = agreement.Сотрудник1
left join stack.[Сотрудники] as staff2 on staff2.ROW_ID = agreement.Сотрудник2
left join stack.[Сотрудники] as staff3 on staff3.ROW_ID = agreement.Сотрудник3
left join stack.[Сотрудники] as staff4 on staff4.ROW_ID = agreement.Сотрудник4
WHERE 	('2023-04-01' between agreement.[Начало договора] and agreement.[Окончание]) and 
		('2023-04-30' between agreement.[Начало договора] and agreement.[Окончание]) and 
		LEN( agreement.[Номер]) = 10) as orul
where fio3 is not null
group by fio3;