from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
import time
import csv
import json

profile = webdriver.FirefoxProfile(
    '/Users/pranshukohli/Library/Application Support/Firefox/Profiles/x122828s.default-release')
profile.set_preference("dom.webdriver.enabled", False)
profile.set_preference('useAutomationExtension', False)
profile.set_preference("browser.download.manager.skipWinSecurityPolicyChecks", True)
profile.set_preference("browser.safebrowsing.allowOverride", True)
profile.update_preferences()

desired = DesiredCapabilities.FIREFOX
driver = webdriver.Firefox(firefox_profile=profile,
                           desired_capabilities=desired)
                           
try:
    driver.get("https://www.etmoney.com/dashboard/mf")

    time.sleep(5)

    login_button='googleLoginBtn'
    element3 = driver.find_element(By.ID, login_button)
    print(element3.text)
    element3.click()
    time.sleep(5)

    driver.switch_to.window(driver.window_handles[1])

    email_id_button = '/html/body/div/div[1]/div/div/main/div[2]/div/div[1]/div[1]'
    element5 = driver.find_element(By.XPATH, email_id_button)
    element5.click()
    print(element5.text)
    time.sleep(5)

    driver.switch_to.window(driver.window_handles[0])


    time.sleep(5)

    
    view_holdings = '/html/body/div[2]/main/div[2]/div/div[1]/div[1]/div/div/div[1]/div[2]/div[2]/button[1]' 
    element8 = driver.find_element(By.XPATH, view_holdings)
    print(element8.text)
    element8.click()

    time.sleep(5)

    bxp1 = "/html/body/div[2]/main/div[2]/div[2]/div/div/div[2]/div[1]/div[2]/div["
    hxp2 = "]/div[2]/div/div[2]/div[1]/div/span/h3"
    vdxp2 = "]/div[2]/div/div[2]/div[2]/div/div[4]/div"

    mfns = []
    for mfn in range(1,14,1):
        _mfn_name = driver.find_element(By.XPATH, bxp1 + str(mfn) + hxp2).text
        print("mfname: " + _mfn_name)
        # _mfn_invested = driver.find_element(By.XPATH, bxp1 + str(mfn) + ixp2).text
        # _mfn_gain = driver.find_element(By.XPATH, bxp1 + str(mfn) + cxp2).text
        # _mfn_xirr = driver.find_element(By.XPATH, bxp1 + str(mfn) + xxp2).text
    
        ele = driver.find_element(By.XPATH, bxp1 + str(mfn) + vdxp2)
        driver.execute_script("arguments[0].click();", ele)
        time.sleep(2)
        bxp11 = "/html/body/div[2]/main/div[2]/div[2]/div/div/div[1]/div[1]/div/div/div[2]/"
        abxp12 ="span"
        ixp12 = "div/div[1]/span[2]"
        cxp12 = "div/div[2]/span[2]"
        txp12 = "div/div[3]/span[2]"
        xxp12 = "div/div[4]/span[2]"
        cnxp12 = "/html/body/div[2]/main/div[2]/div[2]/div/div/div[1]/div[1]/div/div/div[4]/div/div[1]/span"
        anxp12 = "/html/body/div[2]/main/div[2]/div[2]/div/div/div[1]/div[1]/div/div/div[4]/div/div[2]/span"
        uxp12 =  "/html/body/div[2]/main/div[2]/div[2]/div/div/div[1]/div[1]/div/div/div[4]/div/div[3]/span"
        fnxp12 = "/html/body/div[2]/main/div[2]/div[2]/div/div/div[1]/div[3]/div[2]/div/h3"
        

        ele2 = driver.find_element(By.XPATH, bxp11+abxp12)
        driver.execute_script("arguments[0].click();", ele2)
        _mfn_invested = driver.find_element(By.XPATH, bxp11+ixp12).text
        _mfn_ur = driver.find_element(By.XPATH, bxp11+cxp12).text
        _mfn_tr = driver.find_element(By.XPATH, bxp11+txp12).text
        _mfn_xirr = driver.find_element(By.XPATH, bxp11+xxp12).text
        _mfn_cn = driver.find_element(By.XPATH, cnxp12).text
        _mfn_an = driver.find_element(By.XPATH, anxp12).text
        _mfn_units = driver.find_element(By.XPATH, uxp12).text
        time.sleep(2)
        _mfn_fn = driver.find_element(By.XPATH, fnxp12).text

        print("got_values_of_mf")
        time.sleep(5)

        # tvrxp = "/html/body/div[2]/main/div[2]/div[2]/div/div/div[1]/div[5]/div[2]/div[4]/a"
        # ele3 = driver.find_element(By.XPATH, tvrxp)
        # driver.execute_script("arguments[0].click();", ele3)

        # print("transaction_click")
        # time.sleep(10)
        
        # trxpb = "/html/body/div[2]/main/div[2]/div[3]/div/div/div/div[1]/div/ul/div/div/div["
        # investments = []
        # x = 1
        # _inwhile = 1
        # while (_inwhile == 1) :
        #     time.sleep(1)
        #     print("in_while")
        #     dxpn = "]/div[1]/span/p[1]"
        #     vxpn = "]/div[2]"
        #     try:
        #         print(trxpb + str(x) + dxpn)
        #         _mfn_daten = driver.find_element(By.XPATH, trxpb + str(x) + dxpn)
        #         print(_mfn_daten.text)
        #         _mfn_valn = driver.find_element(By.XPATH, trxpb + str(x) + vxpn).text
        #         print("1")
        #         if(x == 11 or x == 13):
        #             driver.execute_script("arguments[0].scrollIntoView();", _mfn_daten)

        #         investments.append(
        #             {
        #                 "date": _mfn_daten.text,
        #                 "value": _mfn_valn
        #             }
        #         )
        #         print("hjh")
        #         x += 1
                
        #     except Exception as e:
        #         print("break"+e)
        #         _inwhile = 0
        #         continue
            
        # print("bn")
        # driver.back()

        print("exit_while")
        mf_json = {
                    "name": _mfn_name,
                    "invested": _mfn_invested.replace("₹","").replace(",",""),
                    "unrealised_returns": _mfn_ur.split(' ')[0].replace("₹","").replace(",",""),
                    "total_returns": _mfn_tr.split(' ')[0].replace("₹","").replace(",",""),
                    "xirr(%)": _mfn_xirr.replace("%","").replace(",",""),
                    "current_nav": _mfn_cn.split(' ')[2].replace("₹","").replace(",",""),
                    "average_nav": _mfn_an.split(' ')[2].replace("₹","").replace(",",""),
                    "units": _mfn_units.split(' ')[1].replace(",",""),
                    "folio_number": _mfn_fn.split(' ')[1],
                    # "investments": investments
                }
        mfns.append(mf_json)
        driver.back()
        time.sleep(5)

    
    print(mfns)
    # headers = ["name","invested","unrealised_returns","total_returns",
    #             "xirr","current_nav","average_nav"]
    headers = mfns[0].keys()

    with open('file.csv', 'w') as f:
        writer = csv.DictWriter(f, fieldnames=headers)
        writer.writeheader()
        writer.writerows(mfns)
    
except Exception as e:
    print(f"Error: {e}")

finally:
    # Quit the browser
    driver.quit()
